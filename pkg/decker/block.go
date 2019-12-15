package decker

import (
  "github.com/stevenaldinger/decker/pkg/hcl"
	"github.com/stevenaldinger/decker/pkg/paths"

	hashicorpHCL "github.com/hashicorp/hcl2/hcl"
)

type Block struct {
  ResourceName     string
  PluginName       string
  PluginConfigPath string
  ForEach          bool
  HCLConfig        *hcl.PluginConfig
  PluginContent    *hashicorpHCL.BodyContent
  Plugins          []Plugin
}

// NewBlock initializes a block.
func NewBlock(block *hashicorpHCL.Block) Block {
  b := Block{
    PluginName: block.Labels[0],
    ResourceName: block.Labels[1],
  }

  b.PluginConfigPath = paths.GetPluginHCLFilePath(b.PluginName)

  pluginAttrs := hcl.GetPluginAttributes(block)

  b.ForEach = contains(pluginAttrs, "for_each")

  b.HCLConfig, b.PluginContent = hcl.GetPluginContent(b.ForEach, block, b.PluginConfigPath)

  return b
}
