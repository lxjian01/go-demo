package services

import (
	"go-demo/httpd/models"
	"go-demo/httpd/utils"
	"go-demo/internal/postgresclient"
)

func AddUser(m *models.User) (int, error) {
	err := postgresclient.DB().Table("user").Create(m).Error
	if err != nil {
		return 0, err
	}
	return m.Id, nil
}

func GetUserPage(pageIndex int, pageSize int, name string) (*utils.PageData, error) {
	dataList := make([]models.User, 0)
	tx := postgresclient.DB().Table("user")
	if name != "" {
		likeStr := "%" + name + "%"
		tx.Where("name like ?", likeStr)
	}
	pageData, err := utils.GetPageData(tx, pageIndex, pageSize, &dataList)
	if err != nil {
		return nil, err
	}
	pageData.Data = &dataList
	return pageData, nil
}
