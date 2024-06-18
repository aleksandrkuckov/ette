import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func createLeaderboard(ctx context.Context, w io.Writer, db string) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database: db,
		Statements: []string{
			`CREATE TABLE Performances (
				PlayerId		INT64 NOT NULL,
				Score		INT64 NOT NULL,
				Timestamp	TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
			) PRIMARY KEY (PlayerId)`,
			`CREATE INDEX PerformancesByScore ON Performances(Score)`,
		},
	})
	if err != nil {
		return err
	}
	if err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created Performances table in database [%s]\n", db)
	return nil
}
  
