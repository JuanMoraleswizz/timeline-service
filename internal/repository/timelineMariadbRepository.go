package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"uala.com/timeline-service/internal/database"
	"uala.com/timeline-service/internal/entity"
)

type TimelineMariadbRepository struct {
	db    database.Database
	redis database.RedisDatabase
}

func NewTimelineMariadbRepository(db database.Database, redis database.RedisDatabase) TimeLineRepository {
	return &TimelineMariadbRepository{
		db:    db,
		redis: redis,
	}
}

func (r *TimelineMariadbRepository) GetTimeline(user int32, forceSync bool) ([]entity.Timeline, error) {
	var timeline []entity.Timeline
	var result []entity.Timeline
	ctxa := context.Background()

	if forceSync {
		r.redis.GetRedis().Del(ctxa, fmt.Sprintf("user:%d", user))
	}
	value, errCache := r.redis.GetRedis().Get(ctxa, fmt.Sprintf("user:%d", user)).Result()

	if errCache != nil {
		fmt.Printf("Error getting cache:%s", errCache.Error())
		fmt.Println("")
	}

	if value != "" {
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			fmt.Printf("Error unmarshalling timeline:%s", err.Error())
			fmt.Println("")
		}
		fmt.Println("result:", result)
		return result, nil
	}

	err := r.db.GetDb().Raw(`
		WITH user_tweets AS (
    		SELECT t.id , t.user_id, u.user_name, t.content, t.create_at
    		FROM tweets  t
    		JOIN users u ON t.user_id = u.id
    		WHERE t.user_id = ?
		),
		followed_tweets AS (
    		SELECT t.id, t.user_id, u.user_name , t.content, t.create_at
    		FROM tweets t
    		JOIN follows f ON t.user_id = f.followee_id
    		JOIN users u ON t.user_id = u.id
    		WHERE f.follower_id = ?
		)
	SELECT *
	FROM user_tweets
	UNION
	SELECT *
	FROM followed_tweets
	ORDER BY create_at DESC;
	`, user, user).Scan(&timeline).Error

	if err != nil {
		return timeline, err
	}

	jsonBytes, err := json.Marshal(timeline)
	if err != nil {
		fmt.Printf("Error marshalling timeline:%s", err.Error())
		fmt.Println("")

	}
	r.redis.GetRedis().Set(ctxa, fmt.Sprintf("user:%d", user), jsonBytes, 0)

	return timeline, nil
}
