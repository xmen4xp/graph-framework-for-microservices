package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"
	crd_generator "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/crd-generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

var _ = Describe("Pkg tests", func() {
	var (
		pkgs   map[string]parser.Package
		pkg    parser.Package
		gnsPkg parser.Package
		ok     bool
	)

	BeforeEach(func() {
		_, err := config.LoadConfig("../../example/nexus-sdk.yaml")
		Expect(err).To(Not(HaveOccurred()))

		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel"]
		Expect(ok).To(BeTrue())
		gnsPkg, ok = pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"]
		Expect(ok).To(BeTrue())
	})

	It("should return generated import strings", func() {
		aliasNameMap := make(map[string]string)
		imports := crd_generator.GenerateImports(&pkg, aliasNameMap)

		expectedImports := []string{
			"configtsmtanzuvmwarecomv1 \"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/config.tsm.tanzu.vmware.com/v1\"",
			"nexustsmtanzuvmwarecomv1 \"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/nexus.tsm.tanzu.vmware.com/v1\""}

		Expect(imports).To(Equal(expectedImports))
	})

	It("should check if node is nexus node", func() {
		nodes := pkg.GetNexusNodes()
		Expect(parser.IsNexusNode(nodes[0])).To(BeTrue())
	})

	It("should check if node is nexus singleton node or not", func() {
		nodes := gnsPkg.GetNexusNodes()
		for _, node := range nodes {
			if node.Name.String() == "Gns" {
				Expect(parser.IsNexusNode(node)).To(BeTrue())
				Expect(parser.IsSingletonNode(node)).To(BeFalse())
				Expect(parser.GetStatusField(node)).NotTo(BeNil())
			} else if node.Name.String() == "Dns" {
				Expect(parser.IsNexusNode(node)).To(BeTrue())
				Expect(parser.IsSingletonNode(node)).To(BeTrue())
			}
		}
	})

	It("should get all structs for gns", func() {
		structs := gnsPkg.GetStructs()
		Expect(structs).To(HaveLen(11))
	})

	It("should get all types for gns", func() {
		types := gnsPkg.GetTypes()
		Expect(types).To(HaveLen(9))
	})

	It("should get imports for gns", func() {
		imports := gnsPkg.GetImportStrings()
		Expect(imports).To(HaveLen(7))
	})

	It("should get all nodes for gns", func() {
		nodes := gnsPkg.GetNodes()
		Expect(nodes).To(HaveLen(7))
	})

	It("should get all consts for gns", func() {
		consts := gnsPkg.GetConsts()
		Expect(consts).To(HaveLen(9))
	})

	It("should get child fields", func() {
		nodes := pkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		Expect(childFields).To(HaveLen(1))
	})

	It("should get link fields for gns", func() {
		nodes := gnsPkg.GetNexusNodes()
		linkFields := parser.GetChildFields(nodes[1])
		Expect(linkFields).To(HaveLen(2))
	})

	It("should get spec fields for gns", func() {
		nodes := gnsPkg.GetNexusNodes()
		specFields := parser.GetSpecFields(nodes[1])
		Expect(specFields).To(HaveLen(5))
	})

	It("should get field name", func() {
		nodes := pkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		fieldName, err := parser.GetNodeFieldName(childFields[0])
		Expect(err).NotTo(HaveOccurred())
		Expect(fieldName).To(Equal("Config"))
	})

	It("should get field type", func() {
		nodes := pkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		fieldType := parser.GetFieldType(childFields[0])
		Expect(fieldType).To(Equal("config.Config"))
	})

	It("should check if field is named child", func() {
		nodes := gnsPkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[1])
		isNamed := parser.IsNamedChildOrLink(childFields[0])
		Expect(isNamed).To(BeTrue())
	})

	It("should get field type for MapType", func() {
		nodes := gnsPkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[1])
		fieldType := parser.GetFieldType(childFields[0])
		Expect(fieldType).To(Equal("service_group.SvcGroup"))
	})

	It("should fail if wrong struct tag is given", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		parser.ParseFieldTags("`nexus: \"child\"`")
		Expect(fail).To(BeTrue())
	})
})