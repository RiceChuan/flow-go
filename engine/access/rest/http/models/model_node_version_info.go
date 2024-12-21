/*
 * Access API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

type NodeVersionInfo struct {
	Semver               string           `json:"semver"`
	Commit               string           `json:"commit"`
	SporkId              string           `json:"spork_id"`
	ProtocolVersion      string           `json:"protocol_version"`
	SporkRootBlockHeight string           `json:"spork_root_block_height"`
	NodeRootBlockHeight  string           `json:"node_root_block_height"`
	CompatibleRange      *CompatibleRange `json:"compatible_range,omitempty"`
}