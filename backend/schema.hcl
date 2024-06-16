table "batch_history" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = false
    type = text
  }
  column "started_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "finished_at" {
    null = true
    type = timestamptz
  }
  column "status" {
    null    = false
    type    = text
    default = "working"
  }
  column "options" {
    null = true
    type = json
  }
  primary_key {
    columns = [column.id]
  }
}
table "contests" {
  schema = schema.public
  column "contest_id" {
    null = false
    type = text
  }
  column "start_epoch_second" {
    null = false
    type = bigint
  }
  column "duration_second" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "rate_change" {
    null = false
    type = text
  }
  column "category" {
    null = false
    type = text
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.contest_id]
  }
}
table "difficulties" {
  schema = schema.public
  column "problem_id" {
    null = false
    type = text
  }
  column "slope" {
    null = true
    type = double_precision
  }
  column "intercept" {
    null = true
    type = double_precision
  }
  column "variance" {
    null = true
    type = double_precision
  }
  column "difficulty" {
    null = true
    type = bigint
  }
  column "discrimination" {
    null = true
    type = double_precision
  }
  column "irt_loglikelihood" {
    null = true
    type = double_precision
  }
  column "irt_users" {
    null = true
    type = double_precision
  }
  column "is_experimental" {
    null = true
    type = boolean
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.problem_id]
  }
}
table "languages" {
  schema = schema.public
  column "language" {
    null = false
    type = text
  }
  column "group" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.language]
  }
  index "languages_group_index" {
    columns = [column.group]
  }
}
table "problems" {
  schema = schema.public
  column "problem_id" {
    null = false
    type = text
  }
  column "contest_id" {
    null = false
    type = text
  }
  column "problem_index" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "title" {
    null = false
    type = text
  }
  column "url" {
    null = false
    type = text
  }
  column "html" {
    null = false
    type = text
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.problem_id]
  }
}
table "submission_crawl_history" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "contest_id" {
    null = false
    type = text
  }
  column "started_at" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.id]
  }
  index "submission_crawl_history_contest_id_start_at_index" {
    columns = [column.contest_id, column.started_at]
  }
}
table "submissions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigint
  }
  column "epoch_second" {
    null = false
    type = bigint
  }
  column "problem_id" {
    null = false
    type = text
  }
  column "contest_id" {
    null = true
    type = text
  }
  column "user_id" {
    null = true
    type = text
  }
  column "language" {
    null = true
    type = text
  }
  column "point" {
    null = true
    type = double_precision
  }
  column "length" {
    null = true
    type = integer
  }
  column "result" {
    null = true
    type = text
  }
  column "execution_time" {
    null = true
    type = integer
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  index "submissions_contest_id_index" {
    columns = [column.contest_id]
  }
  index "submissions_execution_time_index" {
    columns = [column.execution_time]
  }
  index "submissions_epoch_second_index" {
    columns = [column.epoch_second]
  }
  index "submissions_language_index" {
    columns = [column.language]
  }
  index "submissions_length_index" {
    columns = [column.length]
  }
  index "submissions_point_index" {
    columns = [column.point]
  }
  index "submissions_problem_id_index" {
    columns = [column.problem_id]
  }
  index "submissions_result_index" {
    columns = [column.result]
  }
  index "submissions_user_id_index" {
    columns = [column.user_id]
  }
  index "submissions_updated_at_index" {
    columns = [column.epoch_second, column.updated_at]
  }
  partition {
    type    = RANGE
    columns = [column.epoch_second]
  }
}
table "users" {
  schema = schema.public
  column "user_id" {
    null = false
    type = text
  }
  column "rating" {
    null = false
    type = integer
  }
  column "highest_rating" {
    null = false
    type = integer
  }
  column "affiliation" {
    null = true
    type = text
  }
  column "birth_year" {
    null = true
    type = integer
  }
  column "country" {
    null = true
    type = text
  }
  column "crown" {
    null = true
    type = text
  }
  column "join_count" {
    null = false
    type = integer
  }
  column "rank" {
    null = false
    type = integer
  }
  column "active_rank" {
    null = true
    type = integer
  }
  column "wins" {
    null = false
    type = integer
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.user_id]
  }
}
schema "public" {
  comment = "standard public schema"
}
