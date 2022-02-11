package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/logger"
	"github.com/rancher/opni-monitoring/pkg/management"
	"github.com/rancher/opni-monitoring/pkg/test"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

//#region Test Setup
var _ = Describe("Management API User/Subject Access Management Tests", Ordered, func() {
	var environment *test.Environment
	var client management.ManagementClient
	var fingerprint string
	var err error
	BeforeAll(func() {
		fmt.Println("Starting test environment")
		environment = &test.Environment{
			TestBin: "../../../testbin/bin",
			Logger:  logger.New().Named("test"),
		}
		Expect(environment.Start()).To(Succeed())
		client = environment.NewManagementClient()
		Expect(json.Unmarshal(test.TestData("fingerprints.json"), &testFingerprints)).To(Succeed())

		token, err := client.CreateBootstrapToken(context.Background(), &management.CreateBootstrapTokenRequest{
			Ttl: durationpb.New(time.Minute),
		})
		Expect(err).NotTo(HaveOccurred())

		certsInfo, err := client.CertsInfo(context.Background(), &emptypb.Empty{})
		Expect(err).NotTo(HaveOccurred())
		fingerprint = certsInfo.Chain[len(certsInfo.Chain)-1].Fingerprint
		Expect(fingerprint).NotTo(BeEmpty())

		port, errC := environment.StartAgent("test-cluster-id", token, []string{fingerprint})
		promAgentPort := environment.StartPrometheus(port)
		Expect(promAgentPort).NotTo(BeZero())
		Consistently(errC).ShouldNot(Receive())
	})

	AfterAll(func() {
		fmt.Println("Stopping test environment")
		Expect(environment.Stop()).To(Succeed())
	})

	//#endregion

	//#region Happy Path Tests

	It("can return a list of all Cluster IDs that a specific User (Subject) can access", func() {
		_, err = client.CreateRole(context.Background(), &core.Role{
			Id:         "test-role",
			ClusterIDs: []string{"test-cluster-id"},
			MatchLabels: &core.LabelSelector{
				MatchLabels: map[string]string{"test-label": "test-value"},
			},
		},
		)
		Expect(err).NotTo(HaveOccurred())

		_, err = client.CreateRoleBinding(context.Background(), &core.RoleBinding{
			Id:       "test-rolebinding",
			RoleId:   "test-role",
			Subjects: []string{"test-subject"},
		})
		Expect(err).NotTo(HaveOccurred())

		accessList, err := client.SubjectAccess(context.Background(), &core.SubjectAccessRequest{
			Subject: "test-subject",
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(accessList.Items).To(HaveLen(1))
		Expect(accessList.Items[0].Id).To(Equal("test-cluster-id"))
	})

	//#endregion

	//#region Edge Case Tests

	It("cannot return a list of Cluster IDs that a specific User (Subject) cannot access", func() {
		token, err := client.CreateBootstrapToken(context.Background(), &management.CreateBootstrapTokenRequest{
			Ttl: durationpb.New(time.Minute),
		})
		Expect(err).NotTo(HaveOccurred())

		_, errC := environment.StartAgent("new-cluster-id-2", token, []string{fingerprint})
		Consistently(errC).ShouldNot(Receive())

		_, err = client.CreateRole(context.Background(), &core.Role{
			Id:         "test-role-2",
			ClusterIDs: []string{"new-cluster-id-2"},
			MatchLabels: &core.LabelSelector{
				MatchLabels: map[string]string{"test-label-2": "test-value-2"},
			},
		},
		)
		Expect(err).NotTo(HaveOccurred())

		_, err = client.CreateRoleBinding(context.Background(), &core.RoleBinding{
			Id:       "test-rolebinding-2",
			RoleId:   "test-role-2",
			Subjects: []string{"test-subject-2"},
		})
		Expect(err).NotTo(HaveOccurred())

		accessList, err := client.SubjectAccess(context.Background(), &core.SubjectAccessRequest{
			Subject: "test-subject-2",
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(accessList.Items).To(HaveLen(1))
		Expect(accessList.Items[0].Id).To(Equal("new-cluster-id-2"))
	})

	//#endregion
})
