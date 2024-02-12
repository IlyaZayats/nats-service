package repositories

//
//import (
//	"context"
//	"fmt"
//	"github.com/IlyaZayats/dynus/internal/entity"
//	"github.com/IlyaZayats/dynus/internal/interfaces"
//	"github.com/golang-module/carbon/v2"
//	"github.com/jackc/pgx/v5"
//	"github.com/jackc/pgx/v5/pgtype"
//	"github.com/jackc/pgx/v5/pgxpool"
//	"github.com/samber/lo"
//)
//
//type PostgresSlugRepository struct {
//	db *pgxpool.Pool
//}
//
//func NewPostgresSlugRepository(db *pgxpool.Pool) (interfaces.SlugRepository, error) {
//	return &PostgresSlugRepository{
//		db: db,
//	}, nil
//}
//
//func (r *PostgresSlugRepository) UpdateUserSlugs(user entity.User, insertSlugs, deleteSlugs []string, ttl map[string]string) error {
//
//	if len(insertSlugs) != 0 {
//		activeSlugs, err := r.GetActiveSlugs(user)
//		if err != nil {
//			return err
//		}
//
//		var s []string
//		for _, v := range insertSlugs {
//			if !lo.Contains(activeSlugs, v) {
//				s = append(s, v)
//			}
//		}
//		if len(s) != 0 {
//			var str string
//			var ttlIns string
//			q := "INSERT INTO Link (user_id, slug_id, ttl) VALUES "
//			for i := 0; i < len(s); i++ {
//				str = fmt.Sprintf("(%d, (SELECT id FROM Slugs WHERE name='%s'), ", user.Id, s[i])
//				ttlIns = fmt.Sprintf("%s),", "NULL")
//				if i == len(s)-1 {
//					ttlIns = fmt.Sprintf("%s)", "NULL")
//				}
//				if ttl[s[i]] != "" {
//					ttlIns = fmt.Sprintf("'%s'),", ttl[s[i]])
//					if i == len(s)-1 {
//						ttlIns = fmt.Sprintf("'%s')", ttl[s[i]])
//					}
//				}
//				str += ttlIns
//				q += str
//			}
//			_, err := r.db.Exec(context.Background(), q)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	if len(deleteSlugs) != 0 {
//		activeSlugs, err := r.GetActiveSlugs(user)
//		if err != nil {
//			return err
//		}
//
//		var s []string
//		for _, v := range deleteSlugs {
//
//			if lo.Contains(activeSlugs, v) {
//				s = append(s, v)
//			}
//		}
//		if len(s) != 0 {
//			var str, sqlValues string
//			for i := 0; i < len(s); i++ {
//				str = fmt.Sprintf("(%d, (SELECT id FROM Slugs WHERE name='%s')), ", user.Id, s[i])
//				if i == len(s)-1 {
//					str = fmt.Sprintf("(%d, (SELECT id FROM Slugs WHERE name='%s'))", user.Id, s[i])
//				}
//				sqlValues += str
//			}
//			q := fmt.Sprintf("UPDATE Link AS t SET is_valid=False FROM (VALUES %s) AS c(user_id, slug_id) WHERE t.user_id=c.user_id AND t.slug_id=c.slug_id", sqlValues)
//			_, err := r.db.Exec(context.Background(), q)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func (r *PostgresSlugRepository) GetActiveSlugs(user entity.User) ([]string, error) {
//	var activeSlugs []string
//	q := "SELECT Slugs.name FROM Slugs JOIN Link ON Link.slug_id=Slugs.id WHERE Link.is_valid=True AND Link.user_id=$1"
//	rows, err := r.db.Query(context.Background(), q, user.Id)
//	if err != nil && err.Error() != "no rows in result set" {
//		return activeSlugs, err
//	}
//	activeSlugs, err = parseRowsToSlice(rows)
//	return activeSlugs, nil
//}
//
//func parseRowsToSlice(rows pgx.Rows) ([]string, error) {
//	var slice []string
//	defer rows.Close()
//	for rows.Next() {
//		var slug string
//		if err := rows.Scan(&slug); err != nil {
//			return slice, err
//		}
//		slice = append(slice, slug)
//	}
//	return slice, nil
//}
//
//func (r *PostgresSlugRepository) DeleteSlug(slug entity.Slug) error {
//	q := "DELETE FROM Slugs WHERE name=$1"
//	_, err := r.db.Exec(context.Background(), q, slug.Name)
//	return err
//}
//
//func (r *PostgresSlugRepository) InsertSlug(slug entity.Slug) error {
//	q := "INSERT INTO Slugs (name, chance) VALUES ($1, $2)"
//	if _, err := r.db.Exec(context.Background(), q, slug.Name, slug.Chance); err != nil {
//		return err
//	}
//	if slug.Chance > 0 {
//		sql := fmt.Sprintf("INSERT INTO Link (user_id, slug_id) SELECT t.user_id as user_id, (SELECT id from Slugs WHERE name='%s') as slug_id FROM (SELECT id as user_id FROM Users ORDER BY random() LIMIT (SELECT (count(*) * %f) AS user FROM Users) ) as t", slug.Name, slug.Chance)
//		if _, err := r.db.Exec(context.Background(), sql); err != nil {
//			return err
//		}
//		fmt.Println(sql)
//	}
//
//	return nil
//}
//
//func (r *PostgresSlugRepository) CleanupByTTL() error {
//	q := "UPDATE Link SET is_valid=False WHERE created_at+ttl<now() AND is_valid=True"
//	_, err := r.db.Exec(context.Background(), q)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func parseRowsToHistory(rows pgx.Rows, oper string) ([]entity.History, error) {
//	var history []entity.History
//	defer rows.Close()
//	for rows.Next() {
//		var userId, slugName string
//		var dateTime pgtype.Timestamp
//		if err := rows.Scan(&userId, &slugName, &dateTime); err != nil {
//			return history, err
//		}
//		history = append(history, entity.History{UserId: userId, Slug: slugName, Operation: oper, Timestamp: dateTime.Time.String()})
//	}
//	return history, nil
//}
//
//func (r *PostgresSlugRepository) GetHistory(date string) ([]entity.History, error) {
//	var history []entity.History
//	date += "-01"
//	dateCarbon := carbon.Parse(date)
//	upper := dateCarbon.AddMonths(1).ToDateTimeString()
//	lower := dateCarbon.AddMonths(-1).ToDateTimeString()
//	q := "SELECT Link.user_id, Slugs.name, Link.created_at FROM Link " +
//		"JOIN Slugs ON Slugs.id = Link.slug_id " +
//		"WHERE Link.created_at>$1 AND Link.created_at<$2"
//	rows, err := r.db.Query(context.Background(), q, lower, upper)
//	if err != nil {
//		return history, err
//	}
//	history1, err := parseRowsToHistory(rows, "insert")
//	if err != nil {
//		return history, err
//	}
//	q = "SELECT Link.user_id, Slugs.name, Link.updated_at FROM Link " +
//		"JOIN Slugs ON Slugs.id = Link.slug_id " +
//		"WHERE Link.is_valid=False AND Link.updated_at>$1 AND Link.updated_at<$2"
//	rows, err = r.db.Query(context.Background(), q, lower, upper)
//	if err != nil {
//		return history, err
//	}
//	history2, err := parseRowsToHistory(rows, "delete")
//	if err != nil {
//		return history, err
//	}
//	history = append(history1, history2...)
//	return history, nil
//}
