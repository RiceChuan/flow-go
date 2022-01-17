/*
 * Access API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package generated

type AccountPublicKey struct {
	Index string `json:"index"`

	PublicKey string `json:"public_key"`

	SigningAlgorithm *SigningAlgorithm `json:"signing_algorithm"`

	HashingAlgorithm *HashingAlgorithm `json:"hashing_algorithm"`

	SequenceNumber string `json:"sequence_number"`

	Weight string `json:"weight"`

	Revoked bool `json:"revoked"`
}