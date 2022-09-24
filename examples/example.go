package examples

import "github.com/google/uuid"

//Repository type to be mocked
type Repository interface {
	GetAssetStatus(id uuid.UUID) (string, error)
	GetAllAssets() ([]string, error)
	UpsertAsset(id uuid.UUID, status string) error
	DeleteAsset(id uuid.UUID) error
}

//Demo repository consumer
type Demo struct {
	repo Repository
}

func (d Demo) CloseAsset(id uuid.UUID) error {
	return d.repo.UpsertAsset(id, "closed")
}
func (d Demo) IsAssetClosed(id uuid.UUID) (bool, error) {
	status, err := d.repo.GetAssetStatus(id)
	if err != nil {
		return false, err
	}
	return status == "closed", nil
}
func (d Demo) GetNumberOfClosedAssets() (int, error) {
	stats, err := d.repo.GetAllAssets()
	if err != nil {
		return 0, err
	}
	sum := 0
	for _, s := range stats {
		if s == "closed" {
			sum++
		}
	}

	return sum, nil
}
func (d Demo) DeleteIfClosed(id uuid.UUID) error {
	status, err := d.repo.GetAssetStatus(id)
	if err != nil {
		return err
	}

	if status == "closed" {
		return d.repo.DeleteAsset(id)
	}

	return nil
}
