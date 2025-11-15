package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"strconv"
	"ostrichdb-go/src/lib"
	sdk "ostrichdb-go/src/sdk"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	STR = "STRING"
	INT = "INTEGER"
	INT_ARR = "[]INTEGER"
	STR_ARR = "[]STRING"
)

var currentUser User

type Post struct {
	id uint64
	Author *User
	Content string
	NumOfLikes uint64
	NumOfComments uint64
	WhoLikedPost []string
	WhoCommented []string
	Comments []Post
}


var postNames = []string {"author","content","numOfLikes","numOfComments","whoLikedPost","whoCommented"}

type User struct {
	Handle string
	Following []string
	Followers []string
	//can add more here if needed
}

func NewPostBuilder(user *User, content string) Post{
	return Post{
		Author: user,
		Content: content,
		NumOfLikes:  0,
		NumOfComments:  0,
		WhoLikedPost: nil,
		WhoCommented: nil,
	}
}

func NewUserBuilder(name string) User {
	return User{
		Handle: name,
		Following: nil,
		Followers: nil,
	}
}

// "FEED"
func HandleFeed(pageNum int) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "FEED %v\n", pageNum)
	return bufio.NewReader(conn).ReadString('\n')
}

// "POST"
//Upon creation of a new post, create a new recordbuilder and the record itself and append into the "post" cluster
func HandlePost(proj lib.Project, col lib.Collection, p Post)  error {
	newPost := sdk.NewClusterBuilder(&proj, &col, create_post_name(&col))
	sdk.CreateCluster(newPost)
	var record *lib.Record
	for _ , postName:= range postNames{
		switch(postName){
			case "author":
				record = sdk.NewRecordBuilder(&proj, &col, newPost, postName, STR, p.Author.Handle)
				break
			case "content":
				record = sdk.NewRecordBuilder(&proj, &col, newPost, postName, STR, p.Content)
				break
			case "numOfLikes":
				record = sdk.NewRecordBuilder(&proj, &col, newPost, postName, INT, "0")
				break
			case "numOfComments":
				record = sdk.NewRecordBuilder(&proj, &col, newPost, postName, INT, "0")
				break
			case "whoLikedPost":
				record = sdk.NewRecordBuilder(&proj, &col, newPost, postName, STR_ARR, fmt.Sprintf("%v", p.WhoLikedPost))
				break
			case "whoCommented":
				record = sdk.NewRecordBuilder(&proj, &col, newPost, postName, STR_ARR, fmt.Sprintf("%v", p.WhoCommented))
		}

		sdk.CreateRecord(record)
	}

	return nil
}

// "FETCH"
func HandleFetch(post string) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "FETCH %v\n", post)
	return bufio.NewReader(conn).ReadString('\n')
}

// "COMMENT"
func HandleComment(post, content string) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "COMMENT %v %v\n", post, content)
	return bufio.NewReader(conn).ReadString('\n')
}

// "LIKE"
func HandleLike(post string) (string, error) {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	fmt.Fprintf(conn, "LIKE %v\n", post)
	return bufio.NewReader(conn).ReadString('\n')
}

type model struct {
	cursor int
	width  int
	height int
	text   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{
		text: "",
	}
}

func (m model) View() string {
	s := fmt.Sprintf("cursor %v\n", m.cursor)
	s += fmt.Sprintf("1. cobbcoding\nThis is a post.\n\n")
	s += fmt.Sprintf("2. stam\nThis is a post 2.\n\n")
	for range m.height - 8 {
		s += "\n"
	}
	s += fmt.Sprintf(":%v|", m.text)
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			if m.cursor > 0 {
				m.cursor--
				m.text = m.text[:m.cursor]
			}
			break
		case "enter":
			m.text = ""
			m.cursor = 0
			break
		default:
			m.text += msg.(tea.KeyMsg).String()
			m.cursor += 1
		}
	case tea.WindowSizeMsg:
		m.width = msg.(tea.WindowSizeMsg).Width
		m.height = msg.(tea.WindowSizeMsg).Height
		break
	}
	return m, nil
}

func Run() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal("Could not start TUI")
	}
}



func parsePostId(id_s string) *Post {
	ids := strings.Split(id_s, "-")
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		fmt.Println("Could not convert to integer: ", ids[0])
		return nil
	}
	if id >= len(posts) {
		fmt.Printf("Post %v does not exist\n", id)
		return nil
	}
	cur_post := &posts[id]
	for id_i := 1; id_i < len(ids); id_i += 1 {
		id := ids[id_i]
		id_val, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Could not convert to integer: ", id)
			return nil
		}
		if id_val >= len(cur_post.Comments) {
			fmt.Printf("Post %v does not exist\n", id_val)
			return nil
		}
		cur_post = &cur_post.Comments[id_val]
	}
	return cur_post
}

var posts []Post

func parseCommand(command string) (string, bool) {
	currentUser:= NewUserBuilder("@Marshall")
	words := strings.Fields(command)
	var res string
	if len(words) < 2 {
		fmt.Println("Not a valid command")
		return "", false
	}
	switch words[0] {
	// example: POST <message>
	case "POST":
		msg := strings.Join(words[1:], " ")
		newPost := NewPostBuilder(&currentUser, msg)
		newPost.id = uint64(len(posts))
		posts = append(posts, newPost)
		res = fmt.Sprintf("%v", newPost)

		break
	// example: FEED <page count>
	case "FEED":
		res = fmt.Sprintf(res, "%v", posts)
		break
	// example: FETCH <post_id-comment_id...>
	case "FETCH":
		if len(words) < 2 {
			fmt.Println("Not enough arguments to command ", words[0])
			return "", false
		}
		post := parsePostId(words[1])
		if post == nil {
			fmt.Println("Could not parse ", words[1])
			return "", false
		}
		res = fmt.Sprintf(res, "%v", post)
		break
	// example: COMMENT <id-comment_id...> <message>
	case "COMMENT":
		if len(words) < 3 {
			fmt.Println("Not enough arguments to command ", words[0])
			return "", false
		}
		msg := strings.Join(words[2:], " ")
		post := parsePostId(words[1])
		if post == nil {
			fmt.Println("Could not parse ", words[1])
			return "", false
		}
		fmt.Println(post)
		comment := Post{id: uint64(len(post.Comments)), Content: msg}
		post.Comments = append(post.Comments, comment)
		res = "Comment posted successfully..."
		break
	// example: LIKE <id-comment_id..>
	case "LIKE":
		post := parsePostId(words[1])
		if post == nil {
			fmt.Println("Could not parse ", words[1])
			return "", false
		}
		post.NumOfLikes += 1
		res = "Post liked successfully..."
		break
	default:
		log.Fatalf("Unknown command: %s\n", words[0])
	}
	return res, true
}