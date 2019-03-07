package documents

import (
	"bytes"

	"github.com/centrifuge/precise-proofs/proofs"
)

// changedField holds the compact property, old and new value of the field that is changed
// if the old is nil, then it is a set operation
// if new is nil, then it is an unset operation
// if both old and new are set, then it is an edit operation
type changedField struct {
	property, old, new []byte
}

// getChangedFields takes two document trees and returns the compact value, old and new value of the fields that are changed in new tree.
// Properties may have been added to the new tree or removed from the new tree.
// In Either case, since the new tree is different from old, that is considered a change.
func getChangedFields(oldTree, newTree *proofs.DocumentTree, lengthSuffix string) (changedFields []changedField) {
	oldProps := oldTree.PropertyOrder()
	newProps := newTree.PropertyOrder()

	props := make(map[string]proofs.Property)
	for _, p := range append(oldProps, newProps...) {
		// we can ignore the length property since any change in slice or map will return in addition or deletion of properties in the new tree
		if p.Text == lengthSuffix {
			continue
		}

		if _, ok := props[p.ReadableName()]; ok {
			continue
		}

		props[p.ReadableName()] = p
	}

	// check each property and append it changed fields if the value is different.
	for k, p := range props {
		_, ol := oldTree.GetLeafByProperty(k)
		_, nl := newTree.GetLeafByProperty(k)

		if ol == nil {
			changedFields = append(changedFields, newChangedField(p, nl, false))
			continue
		}

		if nl == nil {
			changedFields = append(changedFields, newChangedField(p, ol, true))
			continue
		}

		ov := ol.Value
		nv := nl.Value
		if ol.Hashed {
			ov = ol.Hash
			nv = nl.Hash
		}

		if !bytes.Equal(ov, nv) {
			changedFields = append(changedFields, changedField{
				property: p.CompactName(),
				old:      ov,
				new:      nv,
			})
		}
	}

	return changedFields
}

func newChangedField(p proofs.Property, leaf *proofs.LeafNode, old bool) changedField {
	v := leaf.Value
	if leaf.Hashed {
		v = leaf.Hash
	}

	cf := changedField{property: p.CompactName()}
	if old {
		cf.old = v
		return cf
	}

	cf.new = v
	return cf
}