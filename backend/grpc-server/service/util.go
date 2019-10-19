package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/xerrors"
	"google.golang.org/grpc/metadata"

	pb "github.com/adventar/adventar/backend/grpc-server/grpc/adventar/v1"
	"github.com/adventar/adventar/backend/grpc-server/model"
)

func (s *Service) getCurrentUser(ctx context.Context) (*model.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, xerrors.Errorf("Metadata not found")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, xerrors.Errorf("Authorization metadata not found")
	}

	authResult, err := s.verifier.VerifyIDToken(values[0])
	if err != nil {
		return nil, xerrors.Errorf("Failed to verify token: %w", err)
	}

	var user model.User
	err = s.db.QueryRow("select id, name, icon_url from users where auth_provider = ? and auth_uid = ?", authResult.AuthProvider, authResult.AuthUID).Scan(&user.ID, &user.Name, &user.IconURL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, xerrors.Errorf("Failed query to fetch user: %w", err)
	}

	return &user, nil
}

func (s *Service) bindEntryCount(calendars []*pb.Calendar) error {
	ids := []interface{}{}
	interpolations := []string{}

	for _, c := range calendars {
		ids = append(ids, c.Id)
		interpolations = append(interpolations, "?")
	}

	sql := fmt.Sprintf("select calendar_id, count(*) from entries where calendar_id in (%s) group by calendar_id", strings.Join(interpolations, ","))
	rows, err := s.db.Query(sql, ids...)
	if err != nil {
		return xerrors.Errorf("Failed query to fetch calendar: %w", err)
	}

	entryCounts := map[int64]int32{}
	for rows.Next() {
		var cid int64
		var count int32
		if err := rows.Scan(&cid, &count); err != nil {
			return xerrors.Errorf("Failed to scan row: %w", err)
		}
		entryCounts[cid] = count
	}

	for _, c := range calendars {
		c.EntryCount = entryCounts[c.Id]
	}

	return nil
}

func (s *Service) findEntries(cid int64) ([]*pb.Entry, error) {
	rows, err := s.db.Query(`
		select
			e.id,
			e.day,
			e.title,
			e.comment,
			e.url,
			e.image_url,
			u.id,
			u.name,
			u.icon_url
		from entries as e
		inner join users as u on u.id = e.user_id
		where e.calendar_id = ?
		order by e.day
	`, cid)

	if err != nil {
		return nil, xerrors.Errorf("Failed query to fetch entries: %w", err)
	}

	entries := []*pb.Entry{}
	for rows.Next() {
		var e pb.Entry
		var u pb.User
		err := rows.Scan(
			&e.Id,
			&e.Day,
			&e.Title,
			&e.Comment,
			&e.Url,
			&e.ImageUrl,
			&u.Id,
			&u.Name,
			&u.IconUrl,
		)
		if err != nil {
			return nil, xerrors.Errorf("Failed to scan row: %w", err)
		}
		e.Owner = &u
		entries = append(entries, &e)
	}

	return entries, nil
}

type date struct {
	Year  int
	Month int
	Day   int
}

func currentDate() (*date, error) {
	currentDate := os.Getenv("CURRENT_DATE")
	var t time.Time
	var err error
	if currentDate != "" {
		t, err = time.Parse("2006-01-02", currentDate)
		if err != nil {
			return nil, err
		}
	} else {
		t = time.Now()
	}
	return &date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}, nil
}
