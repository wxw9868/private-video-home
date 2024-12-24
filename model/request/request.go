package request

type Common struct {
	ID uint `uri:"id" binding:"required"`
}

type Paginate struct {
	Page int `uri:"page" form:"page" json:"page"`
	Size int `uri:"size" form:"size" json:"size"`
	// PageSize uint `form:"page_size" json:"page_size"`
}

type OrderBy struct {
	Column string `uri:"column" form:"column" json:"column"`
	Order  string `uri:"order" form:"order" json:"order"`
}

type UpdateUser struct {
	Nickname    string `form:"nickname" json:"nickname" binding:"required"`
	Username    string `form:"username" json:"username" binding:"required"`
	Email       string `form:"email" json:"email" binding:"required,email"`
	Mobile      string `form:"mobile" json:"mobile" binding:"required"`
	Designation string `form:"designation" json:"designation"`
	Address     string `form:"address" json:"address"`
	Intro       string `form:"intro" json:"intro"`
}

type CreateActress struct {
	Name         string `form:"name" json:"name" binding:"required"`
	Alias        string `form:"alias" json:"alias"`
	Avatar       string `form:"avatar" json:"avatar" example:"assets/image/avatar/anonymous.png"`
	Birth        string `form:"birth" json:"birth"`
	Measurements string `form:"measurements" json:"measurements"`
	CupSize      string `form:"cup_size" json:"cup_size"`
	DebutDate    string `form:"debut_date" json:"debut_date"`
	StarSign     string `form:"star_sign" json:"star_sign"`
	BloodGroup   string `form:"blood_group" json:"blood_group"`
	Stature      string `form:"stature" json:"stature"`
	Nationality  string `form:"nationality" json:"nationality"`
	Intro        string `form:"intro" json:"intro"`
}

type UpdateActress struct {
	Id uint `json:"id" binding:"required" `
	CreateActress
}

type SearchActress struct {
	Paginate
	OrderBy
	Actress string `uri:"actress"  form:"actress"  json:"actress"`
}

type SearchVideo struct {
	Paginate
	OrderBy
	ActressID int `uri:"actress_id" form:"actress_id" json:"actress_id"`
}

type CreateDanmu struct {
	VideoID uint    `form:"video_id" json:"video_id" binding:"required"`
	Text    string  `form:"text" json:"text" binding:"required"`
	Time    float64 `form:"time" json:"time"`
	Mode    uint8   `form:"mode" json:"mode"`
	Color   string  `form:"color" json:"color" binding:"required"`
	Border  bool    `form:"border" json:"border"`
	Style   string  `form:"style" json:"style"`
}
