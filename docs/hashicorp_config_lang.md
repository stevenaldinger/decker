# General Guide

- https://github.com/hashicorp/hcl2/tree/master/guide

# Basic HCL Parsing

- [Parsing HCL](https://github.com/hashicorp/hcl2/blob/master/guide/go_parsing.rst)


# Expression Evaluation (used for dynamic HCL reading like whats used for allowing plugins to define the inputs that are read from the user-supplied configs)

- [Expression Evaluation](https://github.com/hashicorp/hcl2/blob/master/guide/go_expression_eval.rst)
- [HCL Low Level Decoding Docs](https://github.com/hashicorp/hcl2/blob/master/guide/go_decoding_lowlevel.rst)
- [Helpful code to see AttributeSchema dynamically built](https://github.com/hashicorp/hcl2/blob/master/gohcl/schema.go)
- [converting between gocty and go](https://github.com/zclconf/go-cty/blob/master/docs/gocty.md)
- [cty/gocty Docs](https://godoc.org/github.com/apparentlymart/go-cty/cty/gocty)

Example code:

```go
parser := hclparse.NewParser()

f, diags := parser.ParseHCLFile("./path/to/hcl_config_file.hcl")

config_file_schema := GetConfigFileSchema()
// Example config file schema:
// &hcl.BodySchema{
//  Blocks: []hcl.BlockHeaderSchema{
//    {
//      Type: "terraform",
//    },
//    {
//      Type:       "resource",
//      LabelNames: []string{"plugin_name", "unique_id"},
//    },
//    {
//      Type:       "variable",
//      LabelNames: []string{"name"},
//    },
//  },
// }

content, _ := f.Body.Content(config_file_schema)

for _, block := range content.Blocks {
  hcl_body_schema := GetBodySchema()
  // Example body schema:
  // &hcl.BodySchema{
  //   Attributes: []hcl.AttributeSchema{
  //     {
  //       Name: "host",
  //     },
  //     {
  //       Name: "plugin_enabled",
  //     },
  //   },
  // }

  block_content, block_diags := block.Body.Content(hcl_body_schema)

	for key, _ := range block_content.Attributes {
    // example context, can be nil if no variables are used in config
    ctx := &hcl.EvalContext{
         Variables: map[string]cty.Value{
             "name": cty.StringVal("Ermintrude"),
             "age":  cty.NumberIntVal(32),
         },
    }

    cty_value, expr_diags block_content.Attributes[key].Expr.Value(ctx)

    var decoded_val string

    gocty.FromCtyValue(cty_value, &decoded_val)

    fmt.Println("My decoded value", decoded_val)
  }
}
```

# Decoder HCL references

- https://godoc.org/github.com/hashicorp/hcl2/gohcl#example-EncodeIntoBody
- https://godoc.org/github.com/hashicorp/hcl2/hclparse#Parser.ParseHCLFile
- https://github.com/hashicorp/hcl2/blob/master/guide/go_decoding_gohcl.rst#id1
