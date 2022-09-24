package examples

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/moqer/funcmoq"
	"github.com/stretchr/testify/assert"
)

//NewRepoMock Creates a  mock repository
func NewRepoMock(t *testing.T) *RepoMock {
	return &RepoMock{
		getAssetStatus: funcmoq.New(t),
		getAllAssets:   funcmoq.New(t),
		upsertAsset:    funcmoq.New(t),
		deleteAsset:    funcmoq.New(t),
	}
}

// RepoMock the struct has to implement all the methods of the interface,
// and beside that it needs one FuncMoq per method
type RepoMock struct {
	getAssetStatus *funcmoq.FuncMoq
	getAllAssets   *funcmoq.FuncMoq
	upsertAsset    *funcmoq.FuncMoq
	deleteAsset    *funcmoq.FuncMoq
}

func (m RepoMock) GetAssetStatus(id uuid.UUID) (status string, err error) { // declaring the return variables here makes the moq slimmer
	// for a specific key (that is registered as part of the arrange section of the test)
	// the object will retrieve a specific set of objects
	m.getAssetStatus.For(id).Retrieve(&status, &err)
	return status, err
}

func (m RepoMock) GetAllAssets() (assets []string, err error) {
	m.getAllAssets.For().Retrieve(&assets, &err)
	return assets, err
}

func (m RepoMock) UpsertAsset(id uuid.UUID, status string) (err error) {
	m.upsertAsset.For(id, status).Retrieve(&err)
	return err
}

func (m RepoMock) DeleteAsset(id uuid.UUID) (err error) {
	m.deleteAsset.For(id).Retrieve(&err)
	return err
}

///////////////////////////////////// TESTS /////////////////////////////////////

func TestDemo_CloseAsset_Error(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	//registers the "key" (id, "closed") to return errors.New("Failed upsert")
	e := errors.New("Failed upsert")
	m.upsertAsset.With(id, "closed").Returns(e)
	demo := Demo{repo: m}

	//act
	err := demo.CloseAsset(id)

	//assert
	assert.Equal(t, e, err)
	assert.Equal(t, 1, m.upsertAsset.CallCount)
}

func TestDemo_CloseAsset_Success(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.upsertAsset.With(id, "closed").Returns(nil)
	demo := Demo{repo: m}

	//act
	err := demo.CloseAsset(id)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, 1, m.upsertAsset.CallCount)
}

func TestDemo_IsAssetClosed_Error(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.getAssetStatus.With(id).Returns("", errors.New("problem"))
	demo := Demo{repo: m}

	//act
	_, err := demo.IsAssetClosed(id)

	//assert
	assert.NotNil(t, err)
	assert.True(t, m.getAssetStatus.Called())
}

func TestDemo_IsAssetClosed_True(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.getAssetStatus.With(id).Returns("closed", nil)
	demo := Demo{repo: m}

	//act
	isClosed, err := demo.IsAssetClosed(id)

	//assert
	assert.Nil(t, err)
	assert.True(t, m.getAssetStatus.Called())
	assert.Equal(t, true, isClosed)
}

func TestDemo_IsAssetClosed_False(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.getAssetStatus.With(id).Returns("notclosed", nil)
	demo := Demo{repo: m}

	//act
	isClosed, err := demo.IsAssetClosed(id)

	//assert
	assert.Nil(t, err)
	assert.True(t, m.getAssetStatus.Called())
	assert.Equal(t, false, isClosed)
}

func TestDemo_GetNumberOfClosedAssets_Error(t *testing.T) {
	//arrange
	m := NewRepoMock(t)

	m.getAllAssets.With().Returns(nil, errors.New("problem"))
	demo := Demo{repo: m}

	//act
	_, err := demo.GetNumberOfClosedAssets()

	//assert
	assert.NotNil(t, err)
	assert.True(t, m.getAllAssets.Called())

}

func TestDemo_GetNumberOfClosedAssets_AllClosed(t *testing.T) {
	//arrange
	m := NewRepoMock(t)

	m.getAllAssets.With().Returns([]string{"closed", "closed", "closed", "closed"}, nil)
	demo := Demo{repo: m}

	//act
	no, err := demo.GetNumberOfClosedAssets()

	//assert
	assert.Nil(t, err)
	assert.True(t, m.getAllAssets.Called())
	assert.Equal(t, 4, no)

}

func TestDemo_GetNumberOfClosedAssets_2Closed(t *testing.T) {
	//arrange
	m := NewRepoMock(t)

	m.getAllAssets.With().Returns([]string{"closed", "closed", "open", "clos"}, nil)
	demo := Demo{repo: m}

	//act
	no, err := demo.GetNumberOfClosedAssets()

	//assert
	assert.Nil(t, err)
	assert.True(t, m.getAllAssets.Called())
	assert.Equal(t, 2, no)
}

func TestDemo_DeleteIfClosed_GetAssetStatusError(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.getAssetStatus.With(id).Returns("", errors.New("problem"))
	demo := Demo{repo: m}

	//act
	err := demo.DeleteIfClosed(id)

	//assert
	assert.NotNil(t, err)
	assert.True(t, m.getAssetStatus.Called())
	assert.Equal(t, 0, m.deleteAsset.CallCount)
}

func TestDemo_DeleteIfClosed_DeleteAssetError(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.getAssetStatus.With(id).Returns("closed", nil)
	m.deleteAsset.With(id).Returns(errors.New("problem"))
	demo := Demo{repo: m}

	//act
	err := demo.DeleteIfClosed(id)

	//assert
	assert.NotNil(t, err)
	assert.True(t, m.getAssetStatus.Called())
	assert.Equal(t, 1, m.deleteAsset.CallCount)
}

func TestDemo_DeleteIfClosed_NotClosed(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.getAssetStatus.With(id).Returns("open", nil)
	m.deleteAsset.With(id).Returns(errors.New("problem"))
	demo := Demo{repo: m}

	//act
	err := demo.DeleteIfClosed(id)

	//assert
	assert.Nil(t, err)
	assert.True(t, m.getAssetStatus.Called())
	assert.Equal(t, 0, m.deleteAsset.CallCount)
}

func TestDemo_DeleteIfClosed_Closed(t *testing.T) {
	//arrange
	m := NewRepoMock(t)
	id := uuid.MustParse("fe1f81f2-6172-4f1a-9377-692a022aec88")

	m.getAssetStatus.With(id).Returns("closed", nil)
	m.deleteAsset.With(id).Returns(nil)
	demo := Demo{repo: m}

	//act
	err := demo.DeleteIfClosed(id)

	//assert
	assert.Nil(t, err)
	assert.True(t, m.getAssetStatus.Called())
	assert.Equal(t, 1, m.deleteAsset.CallCount)
}
