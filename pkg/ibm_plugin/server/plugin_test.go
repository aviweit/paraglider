//go:build unit

/*
Copyright 2023 The Invisinets Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ibm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/platform-services-go-sdk/globalsearchv2"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/require"

	fake "github.com/NetSys/invisinets/pkg/fake/controller/rpc"
	sdk "github.com/NetSys/invisinets/pkg/ibm_plugin/sdk"
	"github.com/NetSys/invisinets/pkg/invisinetspb"
	utils "github.com/NetSys/invisinets/pkg/utils"
)

const (
	fakeResGroup  = "invisinets-fake"
	fakeRegion    = "us-east"  // Primary region used for tests
	fakeConRegion = "us-south" // Region used to test VPC connectivity across regions
	fakeZone      = fakeRegion + "-a"
	fakeInstance  = "vm-invisinets-fake"
	fakeImage     = "fake-image"
	fakeVPC       = "invisinets-fake-vpc"
	fakeID        = "12345"
	fakeID2       = "123452"
	fakeRuleID1   = "fake-rule1"
	fakeRuleID2   = "fake-rule2"
	fakeCRN       = "crn:" + fakeID
	fakeCRN2      = "crn:" + fakeID2
	fakeSubnet    = "invisinets-fake-subnet"
	fakeSG        = "invisinets-fake-sg"
	fakeGW        = "invisnets-fake-gw"
	fakeIP        = "10.0.0.2"
	fakeSubnet1   = "10.0.0.0/16"
	fakeSubnet2   = "20.1.1.0/28"
	fakeProfile   = "bx2-2x8"
	invTag        = "inv"

	fakeResourceID = "/ResourceGroupName/" + fakeResGroup + "/Zone/" + fakeZone + "/ResourceID/" + fakeInstance
	fakeNamespace  = "inv-namespace"
	wrongNamespace = "wrong-inv-namespace"
)

var (
	fakeInstancePrototype = vpcv1.InstancePrototypeInstanceByImage{
		Image:   &vpcv1.ImageIdentityByID{ID: core.StringPtr(fakeImage)},
		Zone:    &vpcv1.ZoneIdentityByName{Name: core.StringPtr(fakeZone)},
		Name:    core.StringPtr(fakeInstance),
		Profile: &vpcv1.InstanceProfileIdentityByName{Name: core.StringPtr(fakeProfile)},
	}

	fakePermitList1 = []*invisinetspb.PermitListRule{
		{
			Id:        fakeRuleID1,
			Direction: invisinetspb.Direction_INBOUND,
			SrcPort:   443,
			DstPort:   443,
			Protocol:  6,
			Targets:   []string{"10.0.0.0/18"},
		},
		{
			Id:        fakeRuleID2,
			Direction: invisinetspb.Direction_OUTBOUND,
			SrcPort:   -1,
			DstPort:   -1,
			Protocol:  -1,
			Targets:   []string{"10.0.64.1"},
		},
	}
	fakePermitList2 = []*invisinetspb.PermitListRule{
		{
			Id:        fakeRuleID1,
			Direction: invisinetspb.Direction_INBOUND,
			SrcPort:   443,
			DstPort:   443,
			Protocol:  6,
			Targets:   []string{"20.1.1.5"},
		},
	}
)

// State of the fake IBM server
// It contains only the necessary items needed to test.
type fakeIBMServerState struct {
	VPCs          []*vpcv1.VPC
	Instance      *vpcv1.Instance
	SecurityGroup *vpcv1.SecurityGroup
	subnetVPC     map[string]string // VPC to Subnet CIDR mapping
}

func sendFakeResponse(w http.ResponseWriter, response interface{}) {
	jsonResp, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Failed to marshal")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, "unable to write request: "+err.Error(), http.StatusBadRequest)
	}
}

// Creates a fake VM Instance.
func createFakeInstance() *vpcv1.Instance {
	return &vpcv1.Instance{
		CRN:    core.StringPtr(fakeCRN),
		Name:   core.StringPtr(fakeInstance),
		ID:     core.StringPtr(fakeID),
		Status: core.StringPtr(vpcv1.InstanceStatusRunningConst),
		NetworkInterfaces: []vpcv1.NetworkInterfaceInstanceContextReference{
			{
				PrimaryIP: &vpcv1.ReservedIPReference{Address: core.StringPtr(fakeIP)},
			},
		},
		VPC: &vpcv1.VPCReference{
			ID:  core.StringPtr(fakeID),
			CRN: core.StringPtr(fakeCRN),
		},
	}
}

// Creates a fake VPC. If connectVPC is set to true, creates another VPC which will be used in multi-VPC test
func createFakeVPC(connectVPC bool) []*vpcv1.VPC {
	var vpcs []*vpcv1.VPC
	vpc := vpcv1.VPC{
		CRN:  core.StringPtr(fakeCRN),
		Name: core.StringPtr(fakeVPC),
		ID:   core.StringPtr(fakeID),
	}
	vpcs = append(vpcs, &vpc)
	if connectVPC {
		vpc2 := vpcv1.VPC{
			CRN:  core.StringPtr(fakeCRN2),
			Name: core.StringPtr(fakeVPC),
			ID:   core.StringPtr(fakeID2),
		}
		vpcs = append(vpcs, &vpc2)
	}
	return vpcs
}

// Creates fake security group. If addRules is set to true, it adds fakePermitList1's rules to it.
func createFakeSecurityGroup(addRules bool) *vpcv1.SecurityGroup {
	sg := vpcv1.SecurityGroup{
		CRN:  core.StringPtr(fakeCRN),
		Name: core.StringPtr(fakeSG),
		ID:   core.StringPtr(fakeID),
	}

	if addRules {
		sgRules := []vpcv1.SecurityGroupRuleIntf{
			&vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp{
				ID:        core.StringPtr(fakeRuleID1),
				Direction: core.StringPtr("inbound"),
				Protocol:  core.StringPtr("tcp"),
				PortMin:   core.Int64Ptr(443),
				PortMax:   core.Int64Ptr(443),
				Remote:    &vpcv1.SecurityGroupRuleRemoteCIDR{CIDRBlock: core.StringPtr("10.0.0.0/18")},
			},
			&vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll{
				ID:        core.StringPtr(fakeRuleID2),
				Direction: core.StringPtr("outbound"),
				Protocol:  core.StringPtr("all"),
				Remote:    &vpcv1.SecurityGroupRuleRemoteIP{Address: core.StringPtr("10.0.64.1")},
			},
		}
		sg.Rules = sgRules
	}

	return &sg
}

// getFakeIBMServerHandler returns the handler with a fake implementation of the IBM Cloud SDK that can be mounted on a URL.
func getFakeIBMServerHandler(fakeIBMServerState *fakeIBMServerState) http.HandlerFunc {
	// The handler should be written as minimally as possible to minimize maintenance overhead. Modifying requests (e.g. POST, DELETE)
	// should generally not do anything other than return the operation response. Instead, initialize the fakeIBMServerState as necessary.
	// Keep in mind these unit tests should rely as little as possible on the functionality of this fake server.

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("unsupported request: %s %s", r.Method, path), http.StatusBadRequest)
			return
		}
		switch {
		case path == "/v3/resources/search":
			if r.Method == http.MethodPost { // Search resources like VPC, Security-group, instance, etc
				var req map[string]interface{}
				err := json.Unmarshal(body, &req)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				res := ""
				tags := []string{}
				var searchResult globalsearchv2.ScanResult
				// Parse the query to extract resource "type" and "tags" to search
				for _, tokens := range strings.Split(req["query"].(string), "AND") {
					keys := strings.Split(tokens, ":")
					switch strings.TrimSpace(keys[0]) {
					case "type":
						res = strings.TrimSpace(keys[1])
					case "tags":
						tags = append(tags, strings.TrimSpace(keys[1]))
					}
				}
				// Search the corresponding resource
				switch sdk.TaggedResourceType(res) {
				case sdk.VM:
					if fakeIBMServerState.Instance != nil {
						wrongNS := false // Namespace is passed in as tags in the query
						for _, tag := range tags {
							if tag == wrongNamespace {
								wrongNS = true // If the wrong namespace is passed, the query shouldn't return the instance
							}
						}
						if !wrongNS {
							var resultItem globalsearchv2.ResultItem
							resultItem.CRN = fakeIBMServerState.Instance.CRN
							searchResult.Items = append(searchResult.Items, resultItem)
						}
					}
				case sdk.SG:
					if fakeIBMServerState.SecurityGroup != nil {
						var resultItem globalsearchv2.ResultItem
						resultItem.CRN = fakeIBMServerState.SecurityGroup.CRN
						searchResult.Items = append(searchResult.Items, resultItem)
					}
				case sdk.VPC:
					if fakeIBMServerState.VPCs != nil {
						for i, fakeVPC := range fakeIBMServerState.VPCs {
							var resultItem globalsearchv2.ResultItem
							resultItem.CRN = fakeVPC.CRN
							// First VPC uses the primary region
							resultItem.SetProperty("region", fakeRegion)
							if i == 1 {
								// Second VPC uses the region used to test cross region connectivity
								resultItem.SetProperty("region", fakeConRegion)
							}
							searchResult.Items = append(searchResult.Items, resultItem)
						}
					}
				case sdk.SUBNET:
					if fakeIBMServerState.subnetVPC != nil {
						for _, subnet := range fakeIBMServerState.subnetVPC {
							var resultItem globalsearchv2.ResultItem
							resultItem.SetProperty("region", fakeRegion)
							resultItem.CRN = core.StringPtr("crn:" + subnet)
							searchResult.Items = append(searchResult.Items, resultItem)
						}
					}
				case sdk.GATEWAY:
					// Not Implemented
				}

				sendFakeResponse(w, searchResult)
				return
			}
		case path == "/vpcs":
			if r.Method == http.MethodPost { // Create VPC
				newVPC := vpcv1.VPC{
					CRN:  core.StringPtr(fakeVPC),
					Name: core.StringPtr(fakeVPC),
					ID:   core.StringPtr(fakeID),
				}
				sendFakeResponse(w, newVPC)
				return
			}
		case path == "/v3/tags/attach":
			if r.Method == http.MethodPost { // Attach tag to a resource
				tagResult := globaltaggingv1.TagResults{
					Results: []globaltaggingv1.TagResultsItem{
						{
							IsError: core.BoolPtr(false),
						},
					},
				}
				sendFakeResponse(w, tagResult)
				return
			}
		case path == "/vpcs/"+fakeID+"/address_prefixes":
			if r.Method == http.MethodPost { // Create Address Prefix
				var newVPCPrefix vpcv1.AddressPrefix
				sendFakeResponse(w, newVPCPrefix)
				return
			}
		case path == "/subnets":
			if r.Method == http.MethodPost { // Create Subnet
				subnet := vpcv1.Subnet{
					CRN: core.StringPtr(fakeSubnet),
					ID:  core.StringPtr(fakeID),
				}
				sendFakeResponse(w, subnet)
				return
			}
			if r.Method == http.MethodGet { // Get Subnets in a VPC
				var subnets vpcv1.SubnetCollection
				if fakeIBMServerState.subnetVPC != nil {
					if cidrBlock, ok := fakeIBMServerState.subnetVPC[r.URL.Query().Get("vpc.id")]; ok {
						subnets.Subnets = []vpcv1.Subnet{
							{
								ID:            core.StringPtr(r.URL.Query().Get("vpc.id")),
								Ipv4CIDRBlock: &cidrBlock,
							},
						}
					}
				}
				sendFakeResponse(w, subnets)
				return
			}
		case strings.Contains(path, "/subnets/"):
			if r.Method == http.MethodGet { // Get Subnet Info for a VPC
				index := strings.LastIndex(path, "/")
				if cidrBlock, ok := fakeIBMServerState.subnetVPC[path[index+1:]]; ok {
					subnet := vpcv1.Subnet{Ipv4CIDRBlock: &cidrBlock}
					sendFakeResponse(w, subnet)
					return
				}
			}
		case path == "/keys":
			if r.Method == http.MethodPost { // Create Key
				key := vpcv1.Key{ID: core.StringPtr(fakeID)}
				sendFakeResponse(w, key)
				return
			}
		case path == "/security_groups":
			if r.Method == http.MethodPost { // Create a security group
				sg := vpcv1.SecurityGroup{
					CRN: core.StringPtr(fakeSG),
					ID:  core.StringPtr(fakeID),
				}
				sendFakeResponse(w, sg)
				return
			}
		case path == "/security_groups/"+fakeID+"/rules":
			if r.Method == http.MethodGet { // Get rules of a security group
				if fakeIBMServerState.SecurityGroup == nil {
					http.Error(w, "Security Group not found", http.StatusNotFound)
					return
				}
				sg := vpcv1.SecurityGroupRuleCollection{
					Rules: fakeIBMServerState.SecurityGroup.Rules,
				}
				sendFakeResponse(w, sg)
				return
			}
			if r.Method == http.MethodPost { // Add rules to a security group
				var sg vpcv1.SecurityGroupRuleIntf
				sendFakeResponse(w, sg)
				return
			}
		case strings.Contains(path, "/security_groups/"+fakeID+"/rules/"):
			if r.Method == http.MethodDelete { // Delete a rule
				w.WriteHeader(http.StatusOK)
				return
			}
		case path == "/instances":
			if r.Method == http.MethodPost { // Create an instance
				fakeIBMServerState.Instance = createFakeInstance()
				sendFakeResponse(w, fakeIBMServerState.Instance)
				return
			}
			if r.Method == http.MethodGet { // List instances
				var instanceCol vpcv1.InstanceCollection
				instanceCol.Instances = make([]vpcv1.Instance, 0)
				if fakeIBMServerState.Instance != nil {
					instanceCol.Instances = append(instanceCol.Instances, *fakeIBMServerState.Instance)
				}
				sendFakeResponse(w, instanceCol)
				return
			}
		case path == "/instances/"+fakeID: // Get an instance
			if r.Method == http.MethodGet {
				if fakeIBMServerState.Instance == nil {
					http.Error(w, "Instance not found", http.StatusNotFound)
					return
				}
				sendFakeResponse(w, fakeIBMServerState.Instance)
				return
			}
		case path == "/instances/"+fakeID+"/network_interfaces":
			if r.Method == http.MethodGet { // List an Instance's network interfaces
				if fakeIBMServerState.Instance == nil {
					http.Error(w, "Instance not found", http.StatusNotFound)
					return
				}
				var netIntf vpcv1.NetworkInterfaceUnpaginatedCollection
				netIntf.NetworkInterfaces = []vpcv1.NetworkInterface{
					{
						SecurityGroups: []vpcv1.SecurityGroupReference{
							{
								Name: fakeIBMServerState.SecurityGroup.Name,
								ID:   fakeIBMServerState.SecurityGroup.ID,
							},
						},
					},
				}
				sendFakeResponse(w, netIntf)
				return
			}
		case path == "/transit_gateways":
			if r.Method == http.MethodPost { // Create transit gateway
				gw := transitgatewayapisv1.TransitGateway{
					Name: core.StringPtr(fakeGW),
					ID:   core.StringPtr(fakeID),
				}
				sendFakeResponse(w, gw)
				return
			}
		case path == "/transit_gateways/"+fakeID+"/connections":
			if r.Method == http.MethodPost { // Create a transit gateway connection
				conn := transitgatewayapisv1.TransitGatewayConnectionCust{
					ID:        core.StringPtr(fakeID),
					Name:      core.StringPtr(fakeGW),
					NetworkID: core.StringPtr("vpc"),
				}
				sendFakeResponse(w, conn)
				return
			}
		}
		fmt.Printf("unsupported request: %s %s\n", r.Method, path)
		http.Error(w, fmt.Sprintf("unsupported request: %s %s", r.Method, path), http.StatusBadRequest)
	})
}

// Creates a http test server, and attaches the fake IBM SDK Handler to it
func setup(t *testing.T, fakeIBMServerState *fakeIBMServerState) (fakeServer *httptest.Server, ctx context.Context, fakeClient *sdk.CloudClient) {
	var err error
	fakeServer = httptest.NewServer(getFakeIBMServerHandler(fakeIBMServerState))
	ctx = context.Background()
	fakeClient, err = sdk.FakeIBMCloudClient(fakeServer.URL, fakeID, fakeRegion)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func TestCreateResourceNewVPC(t *testing.T) {
	_, fakeControllerServerAddr, err := fake.SetupFakeOrchestratorRPCServer(utils.IBM)
	if err != nil {
		t.Fatal(err)
	}
	// fakeIBMServerState with no state to have a clean slate resource creation
	fakeIBMServerState := &fakeIBMServerState{}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		orchestratorServerAddr: fakeControllerServerAddr,
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	description, err := json.Marshal(vpcv1.CreateInstanceOptions{InstancePrototype: vpcv1.InstancePrototypeIntf(&fakeInstancePrototype)})
	require.NoError(t, err)

	resource := &invisinetspb.ResourceDescription{Id: fakeResourceID, Description: description, Namespace: fakeNamespace}
	resp, err := s.CreateResource(ctx, resource)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestCreateResourceExistingVPCSubnet(t *testing.T) {
	_, fakeControllerServerAddr, err := fake.SetupFakeOrchestratorRPCServer(utils.IBM)
	if err != nil {
		t.Fatal(err)
	}
	// fakeIBMServerState with an existing VPC, and subnet
	fakeIBMServerState := &fakeIBMServerState{
		VPCs: createFakeVPC(false),
		subnetVPC: map[string]string{
			fakeID: fakeSubnet1,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		orchestratorServerAddr: fakeControllerServerAddr,
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	description, err := json.Marshal(vpcv1.CreateInstanceOptions{InstancePrototype: vpcv1.InstancePrototypeIntf(&fakeInstancePrototype)})
	require.NoError(t, err)

	resource := &invisinetspb.ResourceDescription{Id: fakeResourceID, Description: description, Namespace: fakeNamespace}
	resp, err := s.CreateResource(ctx, resource)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestCreateResourceExistingVPCMissingSubnet(t *testing.T) {
	_, fakeControllerServerAddr, err := fake.SetupFakeOrchestratorRPCServer(utils.IBM)
	if err != nil {
		t.Fatal(err)
	}
	// fakeIBMServerState with an existing VPC, but no subnet
	fakeIBMServerState := &fakeIBMServerState{
		VPCs: createFakeVPC(false),
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		orchestratorServerAddr: fakeControllerServerAddr,
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	description, err := json.Marshal(vpcv1.CreateInstanceOptions{InstancePrototype: vpcv1.InstancePrototypeIntf(&fakeInstancePrototype)})
	require.NoError(t, err)

	resource := &invisinetspb.ResourceDescription{Id: fakeResourceID, Description: description, Namespace: fakeNamespace}
	resp, err := s.CreateResource(ctx, resource)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestGetUsedAddressSpaces(t *testing.T) {
	// fakeIBMServerState with an existing VPC and subnet
	fakeIBMServerState := &fakeIBMServerState{
		VPCs: createFakeVPC(false),
		subnetVPC: map[string]string{
			fakeID: fakeSubnet1,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	deployment := &invisinetspb.GetUsedAddressSpacesRequest{
		Deployments: []*invisinetspb.InvisinetsDeployment{
			{Id: fakeResourceID, Namespace: fakeNamespace},
		},
	}

	resp, err := s.GetUsedAddressSpaces(ctx, deployment)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.ElementsMatch(t, resp.AddressSpaceMappings[0].AddressSpaces, []string{fakeSubnet1})
}

func TestGetUsedAddressSpacesMultipleVPC(t *testing.T) {
	// fakeIBMServerState with two VPCs (and subnets) across regions
	fakeIBMServerState := &fakeIBMServerState{
		VPCs: createFakeVPC(true),
		subnetVPC: map[string]string{
			fakeID:  fakeSubnet1,
			fakeID2: fakeSubnet2,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion):    fakeClient,
			getClientMapKey(fakeResGroup, fakeConRegion): fakeClient,
		}}

	deployment := &invisinetspb.GetUsedAddressSpacesRequest{
		Deployments: []*invisinetspb.InvisinetsDeployment{
			{Id: fakeResourceID, Namespace: fakeNamespace},
		},
	}

	resp, err := s.GetUsedAddressSpaces(ctx, deployment)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.ElementsMatch(t, resp.AddressSpaceMappings[0].AddressSpaces, []string{fakeSubnet1, fakeSubnet2})
}

func TestAddPermitListRules(t *testing.T) {
	// fakeIBMServerState with an existing VPC, subnet, instance and a security group
	fakeIBMServerState := &fakeIBMServerState{
		VPCs:          createFakeVPC(false),
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(false),
		subnetVPC: map[string]string{
			fakeID: fakeSubnet1,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	addRulesRequest := &invisinetspb.AddPermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		Rules:     fakePermitList1,
	}

	resp, err := s.AddPermitListRules(ctx, addRulesRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAddPermitListRulesExisting(t *testing.T) {
	// fakeIBMServerState with an existing VPC, subnet, instance and a security group
	// with existing rules in fakePermitList1
	fakeIBMServerState := &fakeIBMServerState{
		VPCs:          createFakeVPC(false),
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(true),
		subnetVPC: map[string]string{
			fakeID: fakeSubnet1,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	addRulesRequest := &invisinetspb.AddPermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		Rules:     fakePermitList1,
	}

	resp, err := s.AddPermitListRules(ctx, addRulesRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAddPermitListRulesMissingInstance(t *testing.T) {
	// fakeIBMServerState with empty state without any instance
	fakeIBMServerState := &fakeIBMServerState{}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	addRulesRequest := &invisinetspb.AddPermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		Rules:     fakePermitList1,
	}

	resp, err := s.AddPermitListRules(ctx, addRulesRequest)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAddPermitListRulesMissingSecurityGroup(t *testing.T) {
	// fakeIBMServerState with instance, VPC but no security group created
	fakeIBMServerState := &fakeIBMServerState{
		VPCs:     createFakeVPC(false),
		Instance: createFakeInstance(),
		subnetVPC: map[string]string{
			fakeID: fakeSubnet1,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	addRulesRequest := &invisinetspb.AddPermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		Rules:     fakePermitList1,
	}

	resp, err := s.AddPermitListRules(ctx, addRulesRequest)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAddPermitListRulesWrongNamespace(t *testing.T) {
	// fakeIBMServerState with instance, VPC,  security group and multiple subnets initialized.
	fakeIBMServerState := &fakeIBMServerState{
		VPCs:          createFakeVPC(false),
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(false),
		subnetVPC: map[string]string{
			fakeID:  fakeSubnet1,
			fakeID2: fakeSubnet2,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()

	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	addRulesRequest := &invisinetspb.AddPermitListRulesRequest{
		Namespace: wrongNamespace,
		Resource:  fakeResourceID,
		Rules:     fakePermitList1,
	}

	resp, err := s.AddPermitListRules(ctx, addRulesRequest)
	require.Error(t, err)
	require.Nil(t, resp)
}
func TestAddPermitListRulesTransitGateway(t *testing.T) {
	// fakeIBMServerState with two VPCs (and subnets) across regions
	fakeIBMServerState := &fakeIBMServerState{
		VPCs:          createFakeVPC(true),
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(false),
		subnetVPC: map[string]string{
			fakeID:  fakeSubnet1,
			fakeID2: fakeSubnet2,
		},
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion):    fakeClient,
			getClientMapKey(fakeResGroup, fakeConRegion): fakeClient,
		}}

	// fakePermitList2 is added to the permit list which will trigger creation of a link
	// between VPCs across regions, and hence requiriing deployment of a transit gateway.
	addRulesRequest := &invisinetspb.AddPermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		Rules:     fakePermitList2,
	}

	resp, err := s.AddPermitListRules(ctx, addRulesRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestDeletePermitListRules(t *testing.T) {
	// fakeIBMServerState with an instance and security group with rules
	fakeIBMServerState := &fakeIBMServerState{
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(true),
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	deleteRulesRequest := &invisinetspb.DeletePermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		RuleNames: []string{fakePermitList1[0].Id, fakePermitList1[1].Id},
	}

	resp, err := s.DeletePermitListRules(ctx, deleteRulesRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestDeletePermitListRulesMissingInstance(t *testing.T) {
	// fakeIBMServerState without an instance
	fakeIBMServerState := &fakeIBMServerState{}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	// Currently the plugin takes rule ID since names are not supported by IBM Cloud SDK
	deleteRulesRequest := &invisinetspb.DeletePermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		RuleNames: []string{fakePermitList1[0].Id, fakePermitList1[1].Id},
	}

	resp, err := s.DeletePermitListRules(ctx, deleteRulesRequest)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestDeletePermitListRulesWrongNamespace(t *testing.T) {
	// fakeIBMServerState with a VM and security group
	fakeIBMServerState := &fakeIBMServerState{
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(true),
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	deleteRulesRequest := &invisinetspb.DeletePermitListRulesRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
		RuleNames: []string{fakePermitList1[0].Id, fakePermitList1[1].Id},
	}

	resp, err := s.DeletePermitListRules(ctx, deleteRulesRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
func TestGetPermitList(t *testing.T) {
	// fakeIBMServerState with a VM and security group with rules
	fakeIBMServerState := &fakeIBMServerState{
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(true),
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	getRulesRequest := &invisinetspb.GetPermitListRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
	}

	resp, err := s.GetPermitList(ctx, getRulesRequest)
	require.NoError(t, err)
	require.ElementsMatch(t, resp.Rules, fakePermitList1)
}

func TestGetPermitListEmpty(t *testing.T) {
	// fakeIBMServerState with a VM and security group without rules
	fakeIBMServerState := &fakeIBMServerState{
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(false),
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	getRulesRequest := &invisinetspb.GetPermitListRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
	}

	resp, err := s.GetPermitList(ctx, getRulesRequest)
	require.NoError(t, err)
	require.Empty(t, resp.Rules)
}

func TestGetPermitListMissingInstance(t *testing.T) {
	// fakeIBMServerState with no instance
	fakeIBMServerState := &fakeIBMServerState{}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	getRulesRequest := &invisinetspb.GetPermitListRequest{
		Namespace: fakeNamespace,
		Resource:  fakeResourceID,
	}

	resp, err := s.GetPermitList(ctx, getRulesRequest)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestGetPermitListWrongNamespace(t *testing.T) {
	// fakeIBMServerState with a VM and security group with rules
	fakeIBMServerState := &fakeIBMServerState{
		Instance:      createFakeInstance(),
		SecurityGroup: createFakeSecurityGroup(true),
	}
	fakeServer, ctx, fakeClient := setup(t, fakeIBMServerState)
	defer fakeServer.Close()
	s := &ibmPluginServer{
		cloudClient: map[string]*sdk.CloudClient{
			getClientMapKey(fakeResGroup, fakeRegion): fakeClient,
		}}

	getRulesRequest := &invisinetspb.GetPermitListRequest{
		Namespace: wrongNamespace,
		Resource:  fakeResourceID,
	}

	resp, err := s.GetPermitList(ctx, getRulesRequest)
	require.Error(t, err)
	require.Nil(t, resp)
}