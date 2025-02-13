// Copyright 2022 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package packages_test

import (
	"path/filepath"
	"testing"

	"code.gitea.io/gitea/models/db"
	packages_model "code.gitea.io/gitea/models/packages"
	"code.gitea.io/gitea/models/unittest"
	user_model "code.gitea.io/gitea/models/user"

	_ "code.gitea.io/gitea/models"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		GiteaRootPath: filepath.Join("..", ".."),
	})
}

func TestHasOwnerPackages(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})

	p, err := packages_model.TryInsertPackage(db.DefaultContext, &packages_model.Package{
		OwnerID:   owner.ID,
		LowerName: "package",
	})
	assert.NotNil(t, p)
	assert.NoError(t, err)

	// A package without package versions gets automatically cleaned up and should return false
	has, err := packages_model.HasOwnerPackages(db.DefaultContext, owner.ID)
	assert.False(t, has)
	assert.NoError(t, err)

	pv, err := packages_model.GetOrInsertVersion(db.DefaultContext, &packages_model.PackageVersion{
		PackageID:    p.ID,
		LowerVersion: "internal",
		IsInternal:   true,
	})
	assert.NotNil(t, pv)
	assert.NoError(t, err)

	// A package with an internal package version gets automaticaly cleaned up and should return false
	has, err = packages_model.HasOwnerPackages(db.DefaultContext, owner.ID)
	assert.False(t, has)
	assert.NoError(t, err)

	pv, err = packages_model.GetOrInsertVersion(db.DefaultContext, &packages_model.PackageVersion{
		PackageID:    p.ID,
		LowerVersion: "normal",
		IsInternal:   false,
	})
	assert.NotNil(t, pv)
	assert.NoError(t, err)

	// A package with a normal package version should return true
	has, err = packages_model.HasOwnerPackages(db.DefaultContext, owner.ID)
	assert.True(t, has)
	assert.NoError(t, err)
}
