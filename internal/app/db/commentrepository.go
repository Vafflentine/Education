package db

type CommentRepository struct {
	dbc *DBController
}

/*func (repos *CommentRepository) Create(comment *models.Comment)(*models.Comment,error){
	if repos.dbc.db.Create(&comment)==nil{
		return nil,errors.New("can't insert comment")
	}
	comment.Id=int(comment.ID)
	return comment,nil
}

func (repos *CommentRepository) Delete(comment *models.Comment){
	repos.dbc.db.Delete(&comment)
}

func (repos *CommentRepository) Update(comment *models.Comment,fields map[string]interface{})(*models.Comment,error){
	if repos.dbc.db.Model(&comment).Select(fields).Updates(fields)==nil{
		return nil,errors.New("can't update post")
	}
	return comment,nil
}*/
