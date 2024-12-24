package main

import (
	"fmt"
	"log"
	"flag"
	"github.com/Authing/authing-go-sdk/lib/enum"
	"github.com/Authing/authing-go-sdk/lib/management"
	"github.com/Authing/authing-go-sdk/lib/model"
)

func main() {
	poolID := flag.String("pool", "", "User Pool ID")
	secret := flag.String("secret", "", "API Secret")
	host := flag.String("host", "https://core.authing.cn", "Authing Host")
	iterations := flag.Int("iterations", 1, "测试循环的次数")
	flag.Parse()

	if *poolID == "" || *secret == "" || *host == "" {
		log.Fatal("All parameters --pool, --secret, and --host are required.")
	}

	client := management.NewClient(*poolID, *secret, *host)
	custom := true
	page := 1
	limit := 50

	var allUsers []model.User

	for {
		req := model.QueryListRequest{
			Page:           page,
			Limit:          limit,
			SortBy:         enum.SortByCreatedAtAsc,
			WithCustomData: &custom,
		}
		resp, err := client.GetUserList(req)
		if err != nil {
			log.Fatalf("获取用户列表失败: %v", err)
		}
		allUsers = append(allUsers, resp.List...)
		if len(allUsers) >= int(resp.TotalCount) {
			break
		}
		page++
	}

	for i := 0; i < *iterations; i++ {
		for _, U := range allUsers {
			_, err := client.GetUserDepartments(model.GetUserDepartmentsRequest{Id: U.Id})
			if err != nil {
				log.Printf("获取用户部门失败, 用户 ID %s: %v\n", U.Id, err)
				continue
			}
		}
		fmt.Printf("执行 %d of %d 次\n", i+1, *iterations)
	}
}
