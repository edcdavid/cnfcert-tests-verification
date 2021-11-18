package tests

import (
	"fmt"
	"os"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/nethelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("Networking custom namespace, custom deployment,", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file"))
	}
	execute.BeforeAll(func() {

		By("Clean namespace before all tests")
		err = namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.ApiClient)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {

		By("Clean namespace before each test")
		err := namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.ApiClient)
		Expect(err).ToNot(HaveOccurred())

		By("Remove reports from report directory")
		err = globalhelper.RemoveContentsFromReportDir()
		Expect(err).ToNot(HaveOccurred())
	})

	// 45440
	It("3 custom pods on Default network networking-icmpv4-connectivity", func() {

		By("Define deployment")
		deploymentUnderTest := deployment.RedefineWithReplicaNumber(
			deployment.DefineDeployment(
				netparameters.TestNetworkingNameSpace,
				configSuite.General.TestImage,
				netparameters.TestDeploymentLabels),
			3)

		By("Create Deployment on cluster")
		err := nethelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
		Expect(err).ToNot(HaveOccurred())
		err = os.Setenv(globalparameters.PartnerNamespaceEnvVarName, netparameters.TestNetworkingNameSpace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			[]string{netparameters.NetworkingTestSuiteName},
			netparameters.TestCaseDefaultSkipRegEx,
		)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit report")
		junitTestReport, err := globalhelper.OpenJunitTestReport()
		Expect(err).ToNot(HaveOccurred())
		Expect(globalhelper.IsTestCasePassedInJunitReport(
			junitTestReport,
			netparameters.TestCaseDefaultNetworkName)).To(BeTrue())

		By("Verify test case status in claim report file")
		claimReport, err := globalhelper.OpenClaimReport()
		Expect(err).ToNot(HaveOccurred())
		testPassed, err := globalhelper.IsTestCasePassedInClaimReport(
			netparameters.TestCaseDefaultNetworkName,
			*claimReport)
		Expect(err).ToNot(HaveOccurred())
		Expect(testPassed).To(BeTrue())
	})

	// 45441
	It("custom daemonset, 4 custom pods on Default network", func() {

	})

	// 45442
	It("3 custom pods on Default network networking-icmpv4-connectivity fail when one pod is "+
		"disconnected [negative]", func() {

	})

	// 45443
	It("2 custom pods on Default network networking-icmpv4-connectivity fail when there is no ping binary "+
		"[negative]", func() {

	})

	// 45444
	It("2 custom pods on Default network networking-icmpv4-connectivity skip when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment [skip]", func() {

	})

	// 45445
	It("custom daemonset, 4 custom pods on Default network networking-icmpv4-connectivity pass when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment only", func() {

	})

	// 45446
	It("2 custom pods on Default network networking-icmpv4-connectivity skip when there is no ip binary [skip]",
		func() {

		})

})