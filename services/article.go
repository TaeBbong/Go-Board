package services

import (
	"models"
)

type Article struct {
	ID      int
	Title   string
	Content string
	UserID  uint
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"title":   a.Title,
		"content": a.Content,
		"user_id": a.UserID,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"title":   a.Title,
		"content": a.Content,
	})
}

func (a *Article) Get() (*models.Article, error) {
	var article *models.Article

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles []*models.Article
	)

	articles, err := models.GetArticles(a.getMaps())
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	return maps
}
