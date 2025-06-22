package services

import (
	"mflow/models"

	"gorm.io/gorm"
)

type ArticleService struct {
	DB *gorm.DB
}

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{DB: db}
}

func (s *ArticleService) GetAll() ([]models.Article, error) {
	var articles []models.Article
	err := s.DB.Find(&articles).Error
	return articles, err
}

func (s *ArticleService) GetPaginated(offset, limit int) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	err := s.DB.Model(&models.Article{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.DB.
		Limit(limit).
		Offset(offset).
		Order("tanggal DESC").
		Find(&articles).Error

	return articles, int(total), err
}

func (s *ArticleService) GetPaginatedFiltered(offset, limit int, search, kategori, sort string) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := s.DB.Model(&models.Article{})

	if kategori != "" {
		query = query.Where("kategori LIKE ?", "%"+kategori+"%")
	}

	if search != "" {
		query = query.Where("judul LIKE ? OR konten LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	switch sort {
	case "judul_asc":
		query = query.Order("judul ASC")
	case "judul_desc":
		query = query.Order("judul DESC")
	case "tanggal_asc":
		query = query.Order("tanggal ASC")
	case "tanggal_desc":
		fallthrough
	default:
		query = query.Order("tanggal DESC")
	}

	err = query.Limit(limit).Offset(offset).Find(&articles).Error
	return articles, int(total), err
}

func (s *ArticleService) GetByID(id uint) (models.Article, error) {
	var article models.Article
	err := s.DB.First(&article, id).Error
	return article, err
}

func (s *ArticleService) Create(article *models.Article) error {
	return s.DB.Create(article).Error
}

func (s *ArticleService) Update(id uint, data *models.Article) error {
	var article models.Article
	if err := s.DB.First(&article, id).Error; err != nil {
		return err
	}
	return s.DB.Model(&article).Updates(data).Error
}

func (s *ArticleService) Delete(id uint) error {
	return s.DB.Delete(&models.Article{}, id).Error
}
