package main

import (
	"database/sql"
	"fmt"
	"text/template"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//var DB_CONNECTION_STRING = "Daniil:S_aG@$LmDan37@tcp(80.78.253.5:3306)/richshow"
//var APP_IP = "80.78.253.5"
//var APP_PORT = "8080"

//var CLIENT_NUMBER int16

var DB_CONNECTION_STRING = "mysql:@tcp(127.0.0.1:3306)/new_richshow"

var APP_IP = ""
var APP_PORT = "8080"

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "index", nil)
}

// Динамическая страница с анимационными программами

func animation_article_more(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id_str := vars["id"]

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	child_age_range_id, err := strconv.Atoi(vars["child_age"])
	if err != nil {
		log.Fatal(err)
	}

	type Additional_services struct {
		Id_add             int
		Name_extra_service string
		Extra_services_id  string
		Category_id        string
		Child_age_range_id string
	}

	type Photo_url_struct struct {
		Id_add    string
		Photo_url string
	}

	type Animation_card struct {
		Id                                                uint16
		Name, Full_description, Video_url                 string
		Photo_first                                       string
		Price, Duration, Child_count_min, Child_count_max string
		Photo_url_add                                     []Photo_url_struct
		Add_service                                       []Additional_services
	}

	var showCard = Animation_card{} //

	tmpl, err := template.ParseFiles("templates/animation/animation_article_more.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//db, err := sql.Open("mysql", "mysql:@tcp(127.0.0.1:3306)/new_richshow")

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выбока данных
	res, err := db.Query(fmt.Sprintf("SELECT `id`, `name`, `full_description`, `price`, `duration`, `child_count_min`, `child_count_max` FROM `program` WHERE `id` = %d AND `child_age_range_id` = %d", id, child_age_range_id))
	if err != nil {
		panic(err)
	}

	showCard = Animation_card{} // Возможно это лишняя строка

	for res.Next() {
		var card Animation_card
		err = res.Scan(&card.Id, &card.Name, &card.Full_description, &card.Price, &card.Duration, &card.Child_count_min, &card.Child_count_max)
		if err != nil {
			panic(err)
		}

		//Запрос имен фотографий программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		phs, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id = (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var firstUrlPhoto []string

		for phs.Next() {
			var phs_str string
			err = phs.Scan(&phs_str)
			if err != nil {
				panic(err)
			}
			firstUrlPhoto = append(firstUrlPhoto, phs_str)
		}
		//log.Print(firstUrlPhoto)
		card.Photo_first = firstUrlPhoto[0]

		photoAll, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id <> (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var allUrlPhoto []Photo_url_struct

		for photoAll.Next() {
			var phs_str Photo_url_struct
			err = photoAll.Scan(&phs_str.Photo_url)
			if err != nil {
				panic(err)
			}
			phs_str.Id_add = id_str
			allUrlPhoto = append(allUrlPhoto, phs_str)
		}
		card.Photo_url_add = allUrlPhoto

		// // Присваиваем переменной фотку обложки, первый элемент массива
		// card.Photo_first = allUrlPhoto[0]

		// // Через цикл присваиваем массиву остальные фотографии

		// allUrlPhoto[0] = allUrlPhoto[len(allUrlPhoto)-1]
		// card.Photo_url = allUrlPhoto

		//Запрос текстовго списка дополнительных услуг + их имя, сохраняю в отдельную структуру Add_service
		pt, err := db.Query(fmt.Sprintf("SELECT `extra_services_id`, `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		if err != nil {
			panic(err)
		}

		var extra []Additional_services
		for pt.Next() {
			var pt_int Additional_services

			err = pt.Scan(&pt_int.Extra_services_id, &pt_int.Name_extra_service)
			if err != nil {
				panic(err)
			}

			pt_int.Id_add = id

			program_id, err := strconv.Atoi(pt_int.Extra_services_id)
			if err != nil {
				log.Fatal(err)
			}

			row_link, err := db.Query(fmt.Sprintf("SELECT `category_id`, `child_age_range_id` FROM `program` WHERE `id` = %d", program_id))
			if err != nil {
				panic(err)
			}

			for row_link.Next() {
				err = row_link.Scan(&pt_int.Category_id, &pt_int.Child_age_range_id)
				if err != nil {
					panic(err)
				}
				//log.Print(pt_int.Child_age_range_id)
				//log.Print(pt_int.Category_id)
			}

			extra = append(extra, pt_int)
		}

		card.Add_service = extra

		// type Additional_services struct {
		// 	Name_extra_service string
		// 	Extra_services_id  string
		// 	Category_id        string
		// 	Child_age_range_id string
		// }

		// row_extra, err := db.Query(fmt.Sprintf("SELECT `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		// if err != nil {
		// 	panic(err)
		// }

		// for row_extra.Next() {
		// 	var extra_name string
		// 	err = row_extra.Scan(&extra_name)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	card.Add_service.Name_extra_service = append(card.Add_service.Name_extra_service, extra_name)
		// }

		//Запрос из таблицы url видео

		video, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 2", card.Id))
		if err != nil {
			panic(err)
		}

		video_url_dict := make([]string, 0)

		for video.Next() {
			var video_str string
			err = video.Scan(&video_str)
			if err != nil {
				panic(err)
			}
			video_url_dict = append(video_url_dict, video_str)
			card.Video_url = video_url_dict[0]
		}
		showCard = card
		//log.Print(showCard)
	}
	tmpl.ExecuteTemplate(w, "animation_article_more", showCard)
}

func animation_years(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/animation/animation_years.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "animation_years", nil)

}

// Динамическая страница с анимационными программами для определенного возраста

func animation_article(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	child_age_range_id, err := strconv.Atoi(vars["child_age"])
	if err != nil {
		log.Fatal(err)
	}

	type Animation_cards struct {
		Id                                                  uint16
		Child_age_range_id                                  int
		Name, Short_description, Duration, Price, Photo_url string
	}

	var cards = []Animation_cards{}

	tmpl, err := template.ParseFiles("templates/animation/animation_for_year.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	res, err := db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `duration`, `price` FROM `program` WHERE `child_age_range_id` = %d AND `category_id` = 1 AND (`hit_season` = 1)", child_age_range_id))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Animation_cards
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Duration, &card.Price)
		if err != nil {
			panic(err)

		}

		card.Child_age_range_id = child_age_range_id

		//fmt.Println(fmt.Sprintf("%d", card.Id))

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}

		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}

		cards = append(cards, card)
		//fmt.Println(cards)
	}

	res, err = db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `duration`, `price` FROM `program` WHERE `child_age_range_id` = %d AND `category_id` = 1 AND (`hit_season` = 0)", child_age_range_id))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Animation_cards
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Duration, &card.Price)
		if err != nil {
			panic(err)

		}

		card.Child_age_range_id = child_age_range_id

		//fmt.Println(fmt.Sprintf("%d", card.Id))

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}

		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}

		cards = append(cards, card)
		//fmt.Println(cards)
	}
	tmpl.ExecuteTemplate(w, "animation_for_year", cards) //то что пишешь в шаблоне HTML {{ define "animation_for_year" }}

}

//Заготовка для формы обратной связи (Отправка формы обратной связи через Golang)

func post_information(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/redirect/redirect.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	err = r.ParseForm() // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	//CLIENT_NUMBER += 1

	name := r.FormValue("name")
	tel := r.FormValue("tel")
	utm_source := r.FormValue("utm_source")
	utm_medium := r.FormValue("utm_medium")
	utm_campaign := r.FormValue("utm_campaign")
	utm_content := r.FormValue("utm_content")
	utm_term := r.FormValue("utm_term")

	if utm_source == "" {
		utm_source = "  "
		utm_medium = "  "
		utm_campaign = "  "
		utm_content = "  "
		utm_term = "  "
	}

	// user we are authorizing as
	from := "info@rich-show.ru"

	// use we are sending email to
	to := "amocrm@rich-show.ru"

	// server we are authorized to send email through
	host := "smtp.yandex.ru"

	// Create the authentication for the SendMail()
	// using PlainText, but other authentication methods are encouraged auth := smtp.PlainAuth("", from, "s2BozHwB28S0zQyW", host) mail.ru oaachboahitvekqc
	auth := smtp.PlainAuth("", from, "mcpzwmjnqqzfvwsv", host)

	// NOTE: Using the backtick here ` works like a heredoc, which is why all the
	// rest of the lines are forced to the beginning of the line, otherwise the
	// formatting is wrong for the RFC 822 style
	message := fmt.Sprintf(`
	name: %s
	phone: %s
	Utm_source: %s
	Utm_medium: %s
	Utm_campaign: %s
	Utm_content: %s
	Utm_term: %s`, name, tel, utm_source, utm_medium, utm_campaign, utm_content, utm_term)
	//message := fmt.Sprintf("name: %v", name)
	//fmt.Println(message)

	if err := smtp.SendMail(host+":587", auth, from, []string{to}, []byte(message)); err != nil {
		fmt.Println("Error SendMail: ", err)
		//os.Exit(1)
	}

	var card struct {
		Link string
	}

	card.Link = r.Header.Get("Referer")

	tmpl.ExecuteTemplate(w, "redirect", card)

	// http.Redirect(w, r, r.Header.Get("Referer"), 302) //Установить адрес с надписью "С вами свяжется менеджер в ближайшее время" и кнопкой на главную страницу
}

// Форма для заполнения анимационных программ

func create_article(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/create/create_article.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	type Extra_service_links struct {
		Id       uint16
		Name     string
		Category string
	}

	var linksExtraService = []Extra_service_links{}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	row, err := db.Query("SELECT `id`, `name`, `child_age_range_id` FROM `program`")
	if err != nil {
		panic(err)
	}
	defer row.Close()

	for row.Next() {

		var link Extra_service_links
		err = row.Scan(&link.Id, &link.Name, &link.Category)
		if err != nil {
			panic(err)
		}

		linksExtraService = append(linksExtraService, link)
	}

	tmpl.ExecuteTemplate(w, "create_article", linksExtraService)
}

func save_article(w http.ResponseWriter, r *http.Request) {

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	// for key, values := range r.PostForm {
	// 	fmt.Sprintf("%s - %s", key, values)
	// }
	//fmt.Fprintln(w, r.Form)
	err = r.ParseMultipartForm(8388608) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	formdata := r.MultipartForm

	name := formdata.Value["name"][0]
	category_id := formdata.Value["category_id"][0]
	child_age_range_id := formdata.Value["child_age_range_id"][0]
	full_description := formdata.Value["full_description"][0]
	short_description := formdata.Value["short_description"][0]
	video := formdata.Value["video"][0]
	price := formdata.Value["price"][0]
	child_count_min := formdata.Value["child_count_min"][0]
	child_count_max := formdata.Value["child_count_max"][0]
	duration := formdata.Value["duration"][0]
	hit_season := formdata.Value["hit_season"][0]

	//Перевод значений из str в int
	category_id_int, err := strconv.Atoi(category_id)
	if err != nil {
		log.Fatal(err)
	}

	child_age_range_id_int, err := strconv.Atoi(child_age_range_id)
	if err != nil {
		log.Fatal(err)
	}

	// //Установка данных в таблицу основную
	//log.Print(reflect.TypeOf(price))

	// Возможно следует заменить %s на %d в команде sql ниже
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `program` (`name`, `short_description`, `full_description`, `duration`, `price`, `category_id`, `child_age_range_id`, `child_count_min`, `child_count_max`, `hit_season`) VALUES ('%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s')", name, short_description, full_description, duration, price, category_id_int, child_age_range_id_int, child_count_min, child_count_max, hit_season))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	//Сохраняем в переменной последнее значение id в таблице program
	var program_id int
	err = db.QueryRow("SELECT MAX(ID) FROM `program`").Scan(&program_id)
	if err != nil {
		panic(err)
	}

	//Вызов категории анимационной программы для alt фотографий и видео

	row_category, err := db.Query(fmt.Sprintf("SELECT `name` FROM `category` WHERE `id` = %d", category_id_int))
	if err != nil {
		panic(err)
	}

	name_category := make([]string, 0)

	for row_category.Next() {
		var name string
		if err := row_category.Scan(&name); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		name_category = append(name_category, name)
	}

	program_id_str := strconv.Itoa(program_id)

	// Создание папки для программы

	switch {
	case category_id_int == 1:
		err = os.Mkdir("./static/img/for_article/"+program_id_str+"/", 0755)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

	case category_id_int == 2:
		err = os.Mkdir("./static/img/show_program/"+program_id_str+"/", 0755)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

	case category_id_int == 3:
		err = os.Mkdir("./static/img/master_class/"+program_id_str+"/", 0755)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

	case category_id_int == 4:
		err = os.Mkdir("./static/img/quest/"+program_id_str+"/", 0755)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

	case category_id_int == 5:
		err = os.Mkdir("./static/img/add_services/"+program_id_str+"/", 0755)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	}

	files := formdata.File["photo"] // grab the filenames

	for i := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		name_photo := files[i].Filename

		switch {
		case category_id_int == 1:
			out, err := os.Create("./static/img/for_article/" + program_id_str + "/" + files[i].Filename)

			defer out.Close()
			if err != nil {
				fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
				return
			}

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

		case category_id_int == 2:
			out, err := os.Create("./static/img/show_program/" + program_id_str + "/" + files[i].Filename)

			defer out.Close()
			if err != nil {
				fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
				return
			}

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

		case category_id_int == 3:
			out, err := os.Create("./static/img/master_class/" + program_id_str + "/" + files[i].Filename)

			defer out.Close()
			if err != nil {
				fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
				return
			}

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

		case category_id_int == 4:
			out, err := os.Create("./static/img/quest/" + program_id_str + "/" + files[i].Filename)

			defer out.Close()
			if err != nil {
				fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
				return
			}

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

		case category_id_int == 5:
			out, err := os.Create("./static/img/add_services/" + program_id_str + "/" + files[i].Filename)

			defer out.Close()
			if err != nil {
				fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
				return
			}

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
		}

		//Установка данных в таблицу media

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `media` (`url`, `program_id`, `media_type_id`, `name`, `description`) VALUES ('%s', %d, 1, '%s', '%s')", name_photo, program_id, name, name_category[0]))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		// fmt.Fprintf(w, "Files uploaded successfully : ")
		// fmt.Fprintf(w, files[i].Filename+"\n")

	}

	//Список возможных доп услуг

	extra_services_id_post := formdata.Value["extra_services_id"] // grab the names
	// log.Print(post)

	for i := range extra_services_id_post { // loop through the files one by one
		extra_services_id_range := extra_services_id_post[i]

		extra_services_id_range_int, err := strconv.Atoi(extra_services_id_range)
		if err != nil {
			log.Fatal(err)
		}

		name_extra_row, err := db.Query(fmt.Sprintf("SELECT `name` FROM `program` WHERE `id` = %d", extra_services_id_range_int))
		if err != nil {
			panic(err)
		}

		name_extra := make([]string, 0)

		for name_extra_row.Next() {
			var name string
			if err := name_extra_row.Scan(&name); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				log.Fatal(err)
			}
			name_extra = append(name_extra, name)
		}

		//Установка данных в таблицу дополнительных услуг

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `program_extra_services_links` (`program_id`, `name_main_program`, `extra_services_id`, `name_extra_service`) VALUES (%d, '%s', %d, '%s')", program_id, name, extra_services_id_range_int, name_extra[0]))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
	}

	//Вставка видео в таблицу media

	insert_1, err := db.Query(fmt.Sprintf("INSERT INTO `media` (`url`, `program_id`, `media_type_id`, `name`, `description`) VALUES ('%s', %d, 2, '%s', '%s')", video, program_id, name, name_category[0]))
	if err != nil {
		panic(err)
	}
	defer insert_1.Close()

	// Эксперимент по статье присланной Димой, конец

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Страницы с программами (без разделения на возраст)

func show_programs(w http.ResponseWriter, r *http.Request) {

	type Show_cards struct {
		Id                                        uint16
		Name, Short_description, Price, Photo_url string
		Child_age_range_id string
	}

	var cards = []Show_cards{}

	vars := mux.Vars(r)

	child_age_range_id, err := strconv.Atoi(vars["child_age"])
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("templates/show_program/show_program.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	res, err := db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `price`, `child_age_range_id` FROM `program` WHERE `category_id` = 2 AND `child_age_range_id` = %d ORDER BY `sorting`", child_age_range_id))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Show_cards
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Price, &card.Child_age_range_id)
		if err != nil {
			panic(err)
		}

		//fmt.Println(fmt.Sprintf("%d", card.Id))

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}
		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}
		cards = append(cards, card)
	}

	
	tmpl.ExecuteTemplate(w, "show_program", cards) //то что пишешь в шаблоне HTML {{ define "animation_for_year" }}
}

func show_programs_years(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/show_program/show_programs_years.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "show_programs_years", nil)
}

// Страница с карточкой выбранной программы

func show_program_more(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id_str := vars["id"]

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	type Additional_services struct {
		Name_extra_service string
		Extra_services_id  string
		Category_id        string
		Child_age_range_id string
	}

	type Photo_url_struct struct {
		Id_add    string
		Photo_url string
	}

	type Show_card struct {
		Id                                                uint16
		Name, Full_description, Video_url                 string
		Photo_first                                       string
		Price, Duration, Child_count_min, Child_count_max string
		Photo_url_add                                     []Photo_url_struct
		Add_service                                       []Additional_services
	}

	var showCard = Show_card{} //

	tmpl, err := template.ParseFiles("templates/show_program/show_program_more.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выбока данных
	res, err := db.Query(fmt.Sprintf("SELECT `id`, `name`, `full_description`, `price`, `duration`, `child_count_min`, `child_count_max` FROM `program` WHERE `id` = %d", id))
	if err != nil {
		panic(err)
	}

	showCard = Show_card{} // Возможно это лишняя строка

	for res.Next() {
		var card Show_card
		err = res.Scan(&card.Id, &card.Name, &card.Full_description, &card.Price, &card.Duration, &card.Child_count_min, &card.Child_count_max)
		if err != nil {
			panic(err)
		}

		//Запрос имен фотографий программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		phs, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id = (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var firstUrlPhoto []string

		for phs.Next() {
			var phs_str string
			err = phs.Scan(&phs_str)
			if err != nil {
				panic(err)
			}
			firstUrlPhoto = append(firstUrlPhoto, phs_str)
		}
		//log.Print(firstUrlPhoto)
		card.Photo_first = firstUrlPhoto[0]

		photoAll, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id <> (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var allUrlPhoto []Photo_url_struct

		for photoAll.Next() {
			var phs_str Photo_url_struct
			err = photoAll.Scan(&phs_str.Photo_url)
			if err != nil {
				panic(err)
			}
			phs_str.Id_add = id_str
			allUrlPhoto = append(allUrlPhoto, phs_str)
		}
		card.Photo_url_add = allUrlPhoto

		//Запрос текстовго списка дополнительных услуг + их имя, сохраняю в отдельную структуру Add_service
		pt, err := db.Query(fmt.Sprintf("SELECT `extra_services_id`, `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		if err != nil {
			panic(err)
		}

		var extra []Additional_services
		for pt.Next() {
			var pt_int Additional_services

			err = pt.Scan(&pt_int.Extra_services_id, &pt_int.Name_extra_service)
			if err != nil {
				panic(err)
			}
			program_id, err := strconv.Atoi(pt_int.Extra_services_id)
			if err != nil {
				log.Fatal(err)
			}
			//log.Print(12 + program_id)
			row_link, err := db.Query(fmt.Sprintf("SELECT `category_id`, `child_age_range_id` FROM `program` WHERE `id` = %d", program_id))
			if err != nil {
				panic(err)
			}

			for row_link.Next() {
				err = row_link.Scan(&pt_int.Category_id, &pt_int.Child_age_range_id)
				if err != nil {
					panic(err)
				}
				//log.Print(pt_int.Child_age_range_id)
				//fmt.Print(pt_int.Category_id+"\n")
			}

			extra = append(extra, pt_int)
		}

		card.Add_service = extra

		//Запрос из таблицы url видео

		video, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 2", card.Id))
		if err != nil {
			panic(err)
		}

		video_url_dict := make([]string, 0)

		for video.Next() {
			var video_str string
			err = video.Scan(&video_str)
			if err != nil {
				panic(err)
			}
			video_url_dict = append(video_url_dict, video_str)
			card.Video_url = video_url_dict[0]
		}
		showCard = card
		//log.Print(showCard)
	}
	tmpl.ExecuteTemplate(w, "show_program_more", showCard)
}

func master_class_years(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/master_class/master_class_years.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "master_class_years", nil)
}

func master_class(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/master_class/master_class.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/modul_info.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	vars := mux.Vars(r)

	child_age_range_id, err := strconv.Atoi(vars["child_age"])
	if err != nil {
		log.Fatal(err)
	}

	type Masterclass_card struct {
		Id                                        uint16
		Name, Short_description, Price, Photo_url string
		Child_age_range_id string
	}

	type Masterclass_cards_Category struct {
		Category_id       string
		Masterclass_cards []Masterclass_card
	}

	var cards = Masterclass_cards_Category{}

	cards.Category_id = "3"

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	res, err := db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `price`, `child_age_range_id` FROM `program` WHERE `category_id` = 3 AND `child_age_range_id` = %d ORDER BY `sorting`", child_age_range_id))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Masterclass_card
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Price, &card.Child_age_range_id)
		if err != nil {
			panic(err)
		}

		//fmt.Println(fmt.Sprintf("%d", card.Id))

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}

		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}

		cards.Masterclass_cards = append(cards.Masterclass_cards, card)

	}
	

	//log.Print(cards)

	tmpl.ExecuteTemplate(w, "master_class", cards) //то что пишешь в шаблоне HTML {{ define "animation_for_year" }}

}

// Страница с карточкой выбранной программы

func master_class_more(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id_str := vars["id"]

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	type Additional_services struct {
		Name_extra_service string
		Extra_services_id  string
		Category_id        string
		Child_age_range_id string
	}

	type Photo_url_struct struct {
		Id_add    string
		Photo_url string
	}

	type Show_card struct {
		Id                                                uint16
		Name, Full_description, Video_url                 string
		Photo_first                                       string
		Price, Duration, Child_count_min, Child_count_max string
		Photo_url_add                                     []Photo_url_struct
		Add_service                                       []Additional_services
	}

	// var showCard = Show_card{} //

	tmpl, err := template.ParseFiles("templates/master_class/master_class_more.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выбока данных
	res, err := db.Query(fmt.Sprintf("SELECT `id`, `name`, `full_description`, `price`, `duration`, `child_count_min`, `child_count_max` FROM `program` WHERE `id` = %d ORDER BY `sorting`", id))
	if err != nil {
		panic(err)
	}

	var showCard = Show_card{} // Возможно это лишняя строка

	for res.Next() {
		var card Show_card
		err = res.Scan(&card.Id, &card.Name, &card.Full_description, &card.Price, &card.Duration, &card.Child_count_min, &card.Child_count_max)
		if err != nil {
			panic(err)
		}

		//Запрос имен фотографий программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		phs, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id = (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var firstUrlPhoto []string

		for phs.Next() {
			var phs_str string
			err = phs.Scan(&phs_str)
			if err != nil {
				panic(err)
			}
			firstUrlPhoto = append(firstUrlPhoto, phs_str)
		}
		//log.Print(firstUrlPhoto)
		card.Photo_first = firstUrlPhoto[0]

		photoAll, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id <> (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var allUrlPhoto []Photo_url_struct

		for photoAll.Next() {
			var phs_str Photo_url_struct
			err = photoAll.Scan(&phs_str.Photo_url)
			if err != nil {
				panic(err)
			}
			phs_str.Id_add = id_str
			allUrlPhoto = append(allUrlPhoto, phs_str)
		}
		card.Photo_url_add = allUrlPhoto

		//Запрос текстовго списка дополнительных услуг + их имя, сохраняю в отдельную структуру Add_service
		pt, err := db.Query(fmt.Sprintf("SELECT `extra_services_id`, `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		if err != nil {
			panic(err)
		}

		var extra []Additional_services
		for pt.Next() {
			var pt_int Additional_services

			err = pt.Scan(&pt_int.Extra_services_id, &pt_int.Name_extra_service)
			if err != nil {
				panic(err)
			}
			program_id, err := strconv.Atoi(pt_int.Extra_services_id)
			if err != nil {
				log.Fatal(err)
			}
			//log.Print(12 + program_id)
			row_link, err := db.Query(fmt.Sprintf("SELECT `category_id`, `child_age_range_id` FROM `program` WHERE `id` = %d", program_id))
			if err != nil {
				panic(err)
			}

			for row_link.Next() {
				err = row_link.Scan(&pt_int.Category_id, &pt_int.Child_age_range_id)
				if err != nil {
					panic(err)
				}
				//log.Print(pt_int.Child_age_range_id)
				//log.Print(pt_int.Category_id)
			}

			extra = append(extra, pt_int)
		}

		card.Add_service = extra

		//Запрос из таблицы url видео

		video, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 2", card.Id))
		if err != nil {
			panic(err)
		}

		video_url_dict := make([]string, 0)

		for video.Next() {
			var video_str string
			err = video.Scan(&video_str)
			if err != nil {
				panic(err)
			}
			video_url_dict = append(video_url_dict, video_str)
			card.Video_url = video_url_dict[0]
		}
		showCard = card
		//log.Print(showCard)
	}
	tmpl.ExecuteTemplate(w, "master_class_more", showCard)
}

func quest(w http.ResponseWriter, r *http.Request) {

	type Quest_cards struct {
		Id                                                  uint16
		Child_age_range_id                                  int
		Name, Short_description, Duration, Price, Photo_url string
	}

	var cards = []Quest_cards{}

	tmpl, err := template.ParseFiles("templates/quest/quest.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	res, err := db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `duration`, `price` FROM `program` WHERE `category_id` = 4 AND (`hit_season` = 1)"))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Quest_cards
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Duration, &card.Price)
		if err != nil {
			panic(err)

		}

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}

		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}

		cards = append(cards, card)
		//fmt.Println(cards)
	}

	res, err = db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `duration`, `price` FROM `program` WHERE `category_id` = 4 AND (`hit_season` = 0)"))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Quest_cards
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Duration, &card.Price)
		if err != nil {
			panic(err)

		}

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}

		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}

		cards = append(cards, card)
		//fmt.Println(cards)
	}
	tmpl.ExecuteTemplate(w, "quest", cards) //то что пишешь в шаблоне HTML {{ define "animation_for_year" }}
}

func quest_more(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id_str := vars["id"]

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	type Additional_services struct {
		Id_add             int
		Name_extra_service string
		Extra_services_id  string
		Category_id        string
		Child_age_range_id string
	}

	type Photo_url_struct struct {
		Id_add    string
		Photo_url string
	}

	type Quest_card struct {
		Id                                                uint16
		Name, Full_description, Video_url                 string
		Photo_first                                       string
		Price, Duration, Child_count_min, Child_count_max string
		Photo_url_add                                     []Photo_url_struct
		Add_service                                       []Additional_services
	}

	var showCard = Quest_card{} //

	tmpl, err := template.ParseFiles("templates/quest/quest_more.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выбока данных
	res, err := db.Query(fmt.Sprintf("SELECT `id`, `name`, `full_description`, `price`, `duration`, `child_count_min`, `child_count_max` FROM `program` WHERE `id` = %d", id))
	if err != nil {
		panic(err)
	}

	showCard = Quest_card{} // Возможно это лишняя строка

	for res.Next() {
		var card Quest_card
		err = res.Scan(&card.Id, &card.Name, &card.Full_description, &card.Price, &card.Duration, &card.Child_count_min, &card.Child_count_max)
		if err != nil {
			panic(err)
		}

		//Запрос имен фотографий программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		phs, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id = (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var firstUrlPhoto []string

		for phs.Next() {
			var phs_str string
			err = phs.Scan(&phs_str)
			if err != nil {
				panic(err)
			}
			firstUrlPhoto = append(firstUrlPhoto, phs_str)
		}
		//log.Print(firstUrlPhoto)
		card.Photo_first = firstUrlPhoto[0]

		photoAll, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id <> (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var allUrlPhoto []Photo_url_struct

		for photoAll.Next() {
			var phs_str Photo_url_struct
			err = photoAll.Scan(&phs_str.Photo_url)
			if err != nil {
				panic(err)
			}
			phs_str.Id_add = id_str
			allUrlPhoto = append(allUrlPhoto, phs_str)
		}
		card.Photo_url_add = allUrlPhoto

		// // Присваиваем переменной фотку обложки, первый элемент массива
		// card.Photo_first = allUrlPhoto[0]

		// // Через цикл присваиваем массиву остальные фотографии

		// allUrlPhoto[0] = allUrlPhoto[len(allUrlPhoto)-1]
		// card.Photo_url = allUrlPhoto

		//Запрос текстовго списка дополнительных услуг + их имя, сохраняю в отдельную структуру Add_service
		pt, err := db.Query(fmt.Sprintf("SELECT `extra_services_id`, `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		if err != nil {
			panic(err)
		}

		var extra []Additional_services
		for pt.Next() {
			var pt_int Additional_services

			err = pt.Scan(&pt_int.Extra_services_id, &pt_int.Name_extra_service)
			if err != nil {
				panic(err)
			}

			pt_int.Id_add = id

			program_id, err := strconv.Atoi(pt_int.Extra_services_id)
			if err != nil {
				log.Fatal(err)
			}

			row_link, err := db.Query(fmt.Sprintf("SELECT `category_id`, `child_age_range_id` FROM `program` WHERE `id` = %d", program_id))
			if err != nil {
				panic(err)
			}

			for row_link.Next() {
				err = row_link.Scan(&pt_int.Category_id, &pt_int.Child_age_range_id)
				if err != nil {
					panic(err)
				}
				//log.Print(pt_int.Child_age_range_id)
				//log.Print(pt_int.Category_id)
			}

			extra = append(extra, pt_int)
		}

		card.Add_service = extra

		// type Additional_services struct {
		// 	Name_extra_service string
		// 	Extra_services_id  string
		// 	Category_id        string
		// 	Child_age_range_id string
		// }

		// row_extra, err := db.Query(fmt.Sprintf("SELECT `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		// if err != nil {
		// 	panic(err)
		// }

		// for row_extra.Next() {
		// 	var extra_name string
		// 	err = row_extra.Scan(&extra_name)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	card.Add_service.Name_extra_service = append(card.Add_service.Name_extra_service, extra_name)
		// }

		//Запрос из таблицы url видео

		video, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 2", card.Id))
		if err != nil {
			panic(err)
		}

		video_url_dict := make([]string, 0)

		for video.Next() {
			var video_str string
			err = video.Scan(&video_str)
			if err != nil {
				panic(err)
			}
			video_url_dict = append(video_url_dict, video_str)
			card.Video_url = video_url_dict[0]
		}
		showCard = card
		//log.Print(showCard)
	}
	tmpl.ExecuteTemplate(w, "quest_more", showCard)
}

func add_services(w http.ResponseWriter, r *http.Request) {

	type Add_cards struct {
		Id                                                  uint16
		Child_age_range_id                                  int
		Name, Short_description, Duration, Price, Photo_url string
	}

	var cards = []Add_cards{}

	tmpl, err := template.ParseFiles("templates/add_services/add_services.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	res, err := db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `duration`, `price` FROM `program` WHERE `category_id` = 5 AND (`hit_season` = 1)"))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Add_cards
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Duration, &card.Price)
		if err != nil {
			panic(err)

		}

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}

		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}

		cards = append(cards, card)
		//fmt.Println(cards)
	}

	res, err = db.Query(fmt.Sprintf("SELECT `Id`, `name`, `short_description`, `duration`, `price` FROM `program` WHERE `category_id` = 5 AND (`hit_season` = 0)"))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {

		var card Add_cards
		err = res.Scan(&card.Id, &card.Name, &card.Short_description, &card.Duration, &card.Price)
		if err != nil {
			panic(err)

		}

		//Запрос имени фотографии программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		photo, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 LIMIT 1", card.Id))
		if err != nil {
			panic(err)
		}

		for photo.Next() {
			err = photo.Scan(&card.Photo_url)
			if err != nil {
				panic(err)
			}
		}

		cards = append(cards, card)
		//fmt.Println(cards)
	}
	tmpl.ExecuteTemplate(w, "add_services", cards) //то что пишешь в шаблоне HTML {{ define "animation_for_year" }}

}

func add_services_more(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id_str := vars["id"]

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	type Additional_services struct {
		Id_add             int
		Name_extra_service string
		Extra_services_id  string
		Category_id        string
		Child_age_range_id string
	}

	type Photo_url_struct struct {
		Id_add    string
		Photo_url string
	}

	type Add_card struct {
		Id                                                uint16
		Name, Full_description, Video_url                 string
		Photo_first                                       string
		Price, Duration, Child_count_min, Child_count_max string
		Photo_url_add                                     []Photo_url_struct
		Add_service                                       []Additional_services
	}

	var showCard = Add_card{} //

	tmpl, err := template.ParseFiles("templates/add_services/add_services_more.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	//DB_CONNECTION_STRING := "mysql:@tcp(127.0.0.1:3306)/new_richshow"
	db, err := sql.Open("mysql", DB_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выбока данных
	res, err := db.Query(fmt.Sprintf("SELECT `id`, `name`, `full_description`, `price`, `duration`, `child_count_min`, `child_count_max` FROM `program` WHERE `id` = %d", id))
	if err != nil {
		panic(err)
	}

	showCard = Add_card{} // Возможно это лишняя строка

	for res.Next() {
		var card Add_card
		err = res.Scan(&card.Id, &card.Name, &card.Full_description, &card.Price, &card.Duration, &card.Child_count_min, &card.Child_count_max)
		if err != nil {
			panic(err)
		}

		//Запрос имен фотографий программы для карточки и добавляется также в структуру для дальнейшего вывода на странице
		phs, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id = (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var firstUrlPhoto []string

		for phs.Next() {
			var phs_str string
			err = phs.Scan(&phs_str)
			if err != nil {
				panic(err)
			}
			firstUrlPhoto = append(firstUrlPhoto, phs_str)
		}
		//log.Print(firstUrlPhoto)
		card.Photo_first = firstUrlPhoto[0]

		photoAll, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 1 AND id <> (SELECT MIN(id) FROM `media` WHERE `program_id` = %d)", card.Id, card.Id))
		if err != nil {
			panic(err)
		}

		var allUrlPhoto []Photo_url_struct

		for photoAll.Next() {
			var phs_str Photo_url_struct
			err = photoAll.Scan(&phs_str.Photo_url)
			if err != nil {
				panic(err)
			}
			phs_str.Id_add = id_str
			allUrlPhoto = append(allUrlPhoto, phs_str)
		}
		card.Photo_url_add = allUrlPhoto

		// // Присваиваем переменной фотку обложки, первый элемент массива
		// card.Photo_first = allUrlPhoto[0]

		// // Через цикл присваиваем массиву остальные фотографии

		// allUrlPhoto[0] = allUrlPhoto[len(allUrlPhoto)-1]
		// card.Photo_url = allUrlPhoto

		//Запрос текстовго списка дополнительных услуг + их имя, сохраняю в отдельную структуру Add_service
		pt, err := db.Query(fmt.Sprintf("SELECT `extra_services_id`, `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		if err != nil {
			panic(err)
		}

		var extra []Additional_services
		for pt.Next() {
			var pt_int Additional_services

			err = pt.Scan(&pt_int.Extra_services_id, &pt_int.Name_extra_service)
			if err != nil {
				panic(err)
			}

			pt_int.Id_add = id

			program_id, err := strconv.Atoi(pt_int.Extra_services_id)
			if err != nil {
				log.Fatal(err)
			}

			row_link, err := db.Query(fmt.Sprintf("SELECT `category_id`, `child_age_range_id` FROM `program` WHERE `id` = %d", program_id))
			if err != nil {
				panic(err)
			}

			for row_link.Next() {
				err = row_link.Scan(&pt_int.Category_id, &pt_int.Child_age_range_id)
				if err != nil {
					panic(err)
				}
				//log.Print(pt_int.Child_age_range_id)
				//log.Print(pt_int.Category_id)
			}

			extra = append(extra, pt_int)
		}

		card.Add_service = extra

		// type Additional_services struct {
		// 	Name_extra_service string
		// 	Extra_services_id  string
		// 	Category_id        string
		// 	Child_age_range_id string
		// }

		// row_extra, err := db.Query(fmt.Sprintf("SELECT `name_extra_service` FROM `program_extra_services_links` WHERE `program_id` = %d", card.Id))
		// if err != nil {
		// 	panic(err)
		// }

		// for row_extra.Next() {
		// 	var extra_name string
		// 	err = row_extra.Scan(&extra_name)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	card.Add_service.Name_extra_service = append(card.Add_service.Name_extra_service, extra_name)
		// }

		//Запрос из таблицы url видео

		video, err := db.Query(fmt.Sprintf("SELECT `url` FROM `media` WHERE `program_id` = %d AND `media_type_id` = 2", card.Id))
		if err != nil {
			panic(err)
		}

		video_url_dict := make([]string, 0)

		for video.Next() {
			var video_str string
			err = video.Scan(&video_str)
			if err != nil {
				panic(err)
			}
			video_url_dict = append(video_url_dict, video_str)
			card.Video_url = video_url_dict[0]
		}
		showCard = card
		//log.Print(showCard)
	}
	tmpl.ExecuteTemplate(w, "add_services_more", showCard)
}

func about_us(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/about_us/about_us.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "about_us", nil)
}

func privacy_policy(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/privacy_policy/privacy_policy.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "privacy_policy", nil)
}

func page_404(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/page_404/page_404.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "page_404", nil)
}

func ready_holidays_1_3(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ready_holidays/ready_holidays_1_3.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "ready_holidays_1_3", nil)
}

func ready_holidays_4_6(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ready_holidays/ready_holidays_4_6.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "ready_holidays_4_6", nil)
}

func ready_holidays_7_9(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ready_holidays/ready_holidays_7_9.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "ready_holidays_7_9", nil)
}

func ready_holidays_10_14(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ready_holidays/ready_holidays_10_14.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/modul.html", "templates/includes/form_footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "ready_holidays_10_14", nil)
}

func yandex_title(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/page_404/yandex_4f576b017e6c01b9.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "yandex_title", nil)
}


func ready_holidays_years(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ready_holidays/ready_holidays_years.html", "templates/includes/header.html", "templates/includes/footer.html", "templates/includes/form_footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "ready_holidays_years", nil)
}

// func NotFoundHandler(w http.ResponseWriter, r *http.Request, status int) {
// 	log.Print(12)
// 	w.WriteHeader(status)
// 	log.Print(12)
// 	if status == http.StatusNotFound {
// 		fmt.Fprint(w, "custom 404")
// 	}
// }

// в SQL необходимо создать четыре таблицы с названиями animation_year_0_4 animation_year_5_7
//animation_year_8_10 animation_year_11_13
// далее в каждую таблицу добавить данные по списку Id, Title, FullText, PointText, Price, Time, Count1, Count2

func handleRequest() {
	rtr := mux.NewRouter()
	rtr.NotFoundHandler = http.HandlerFunc(page_404)

	rtr.HandleFunc("/", index).Methods("GET")

	//Анимационные программы
	rtr.HandleFunc("/animation_years/", animation_years).Methods("GET")
	rtr.HandleFunc("/animation_years/{child_age}/", animation_article).Methods("GET")
	rtr.HandleFunc("/animation_years/{child_age}/card/{id:[0-9]+}/", animation_article_more).Methods("GET")

	// Шоу программы
	rtr.HandleFunc("/show_programs/", show_programs_years).Methods("GET")
	rtr.HandleFunc("/show_programs/{child_age}/", show_programs).Methods("GET")
	rtr.HandleFunc("/show_programs/{child_age}/card/{id:[0-9]+}/", show_program_more).Methods("GET")

	// Мастер-классы
	rtr.HandleFunc("/master_class/", master_class_years).Methods("GET")
	rtr.HandleFunc("/master_class_years/{child_age}/", master_class).Methods("GET")
	rtr.HandleFunc("/master_class_years/{child_age}/card/{id:[0-9]+}/", master_class_more).Methods("GET")

	// Квесты
	rtr.HandleFunc("/quest/", quest).Methods("GET")
	rtr.HandleFunc("/quest/card/{id:[0-9]+}/", quest_more).Methods("GET")

	// Дополнительные услуги
	rtr.HandleFunc("/add_services/", add_services).Methods("GET")
	rtr.HandleFunc("/add_services/card/{id:[0-9]+}/", add_services_more).Methods("GET")

	// Готовые праздники
	rtr.HandleFunc("/ready_holidays/1/", ready_holidays_1_3).Methods("GET")
	rtr.HandleFunc("/ready_holidays/2/", ready_holidays_4_6).Methods("GET")
	rtr.HandleFunc("/ready_holidays/3/", ready_holidays_7_9).Methods("GET")
	rtr.HandleFunc("/ready_holidays/4/", ready_holidays_10_14).Methods("GET")
	rtr.HandleFunc("/ready_holidays_years/", ready_holidays_years).Methods("GET")

	// Отправка формы на электронную почту
	rtr.HandleFunc("/post_information/", post_information).Methods("POST")

	// Информация о нас
	rtr.HandleFunc("/about_us/", about_us).Methods("GET")

	//Форма для ввода анимационных программ
	rtr.HandleFunc("/admin/create_article/", create_article).Methods("GET")
	rtr.HandleFunc("/admin/save_article/", save_article).Methods("POST")

	//Политика конфиденциальности
	rtr.HandleFunc("/privacy_policy/", privacy_policy).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	//APP_IP := os.Getenv("APP_IP")
	//APP_PORT := os.Getenv("APP_PORT")

	rtr.HandleFunc("/yandex_4f576b017e6c01b9.html", yandex_title).Methods("GET")

	//APP_IP := ""
	//APP_PORT := "8080"

	//fmt.Println(APP_IP + ":" + APP_PORT)
	http.ListenAndServe(APP_IP+":"+APP_PORT, nil)

	//http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
}
