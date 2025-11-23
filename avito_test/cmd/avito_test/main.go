package main

import (
	"avito_test/internal/app"
	"avito_test/internal/config"
	"avito_test/internal/logger"
	"avito_test/internal/storage/postgres"
	"context"
	"fmt"
	"os"
)

func main() {
	config := config.MustLoad()
	log := logger.Setup(os.Stdout, config.LogLevel())

	ctx := context.Background()

	dbc := config.DBConfig
	storagePath := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbc.USER, dbc.PASS, dbc.HOST, dbc.PORT, dbc.DBName)
	storage, err := postgres.New(ctx, storagePath, log)
	if err != nil {
		fmt.Print(err)
	}
	defer storage.Close()
	app := app.New(*config, log, storage)

	err = app.Run()

	if err != nil {
		log.Warn("server running error", err.Error())
	}
	// err = storage.CreateTeam(ctx, &models.Team{
	// 	TeamName: "MyTeam1",
	// 	Members: []*models.TeamMember{
	// 		{
	// 			UserID:   "1",
	// 			Username: "Thomas",
	// 			IsActive: true,
	// 		},
	// 		{
	// 			UserID:   "2",
	// 			Username: "John",
	// 			IsActive: false,
	// 		},
	// 	},
	// })
	// fmt.Print("%s", err)
	// team, err := storage.GetTeamWithMembers(ctx, "MyTeam1")
	// teamJson, err := json.MarshalIndent(team, "-", "  ")
	// fmt.Printf("%v, %v", string(teamJson), err)

	// reviewers, err := storage.CreatePR(ctx, models.NewPullRequest("154", "my first pr", "2"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// revJSON, _ := json.MarshalIndent(reviewers, "-", "  ")
	// fmt.Printf("%v", string(revJSON))

	// pr, err := storage.PRMarkAsMerged(ctx, "154")
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// b, _ := json.MarshalIndent(pr, "--", "  ")
	// fmt.Printf("%v", string(b))

	// pr, err := storage.Reassign(ctx, "154", "1")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// b, _ := json.MarshalIndent(pr, "--", "  ")
	// fmt.Printf("%v", string(b))
	// user, err := storage.GetReview(ctx, "3")
	// b, _ := json.MarshalIndent(user, "--", "  ")
	// fmt.Printf("%v,  %v", string(b), err)
	fmt.Println("Успешное подключение к PostgreSQL через пул соединений!")

}
