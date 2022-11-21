package schema

import (
	"time"

	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"

	"github.com/drakejin/crawler/internal/storage/db/ent/validate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Page holds the schema definition for the Keyword entity.
type Page struct {
	ent.Schema
}

// Fields of the Page.
func (Page) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique(),
		field.String("crawling_version"),

		field.String("domain").
			Annotations(entsql.Annotation{
				Size: 700,
			}).
			Validate(validate.MaxRuneCount(200)).
			Comment("domain www.example.com"),

		field.String("port").
			Annotations(entsql.Annotation{
				Size: 30,
			}).
			Default("80").
			Validate(validate.MaxRuneCount(10)).
			Comment("port number"),

		field.Bool("is_https").
			Default(false).
			Comment("is used tls/ssl layer flag"),

		field.String("url").
			Annotations(entsql.Annotation{
				Size: 750,
			}).
			NotEmpty().
			Validate(validate.MaxRuneCount(1200)).
			Comment("url for only indexing"),

		field.String("path").
			Annotations(entsql.Annotation{
				Size: 1000,
			}).
			Default("").
			Validate(validate.MaxRuneCount(300)).
			Comment("url.path"),

		field.Text("querystring").
			Default("").
			SchemaType(map[string]string{
				dialect.MySQL: "text",
			}).
			Comment("url.querystring"),

		//field.Text("url").
		//	Default("").
		//	SchemaType(map[string]string{
		//		dialect.MySQL: "text",
		//	}).
		//	Comment("this mean url"),

		field.Int64("count_referred").
			Default(0).
			Comment("how many times referred"),

		field.Enum("status").
			Values("ALLOW", "NOTALLOW", "DELETE").
			Default("ALLOW").
			Comment("해당 row는 쓸 수 있는지? 없는지?"),

		field.Time("created_at").
			Default(time.Now().UTC).
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}).
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).
			Comment("first indexed time"),

		field.String("created_by").
			Annotations(entsql.Annotation{
				Size: 300,
			}).
			Validate(validate.MaxRuneCount(100)).
			Comment("first indexed time by which system"),

		field.Time("updated_at").
			Default(time.Now().UTC).
			UpdateDefault(time.Now().UTC).
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
			}).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).
			Comment("modified time"),

		field.String("updated_by").
			Annotations(entsql.Annotation{
				Size: 300,
			}).
			Validate(validate.MaxRuneCount(100)).
			Comment("modified by which system"),

		field.String("title").
			Default("").
			Comment("html title tag"),

		field.String("description").
			Default("").
			Comment("basic meta tags 'description'"),

		field.String("keywords").
			Default("").
			Comment("basic meta tags 'keywords'"),

		field.String("content_language").
			Default("").
			Comment("basic meta tags 'content-language'"),

		// https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/markup 추후 있으면 추가하기
		field.String("twitter_card").
			Default("").
			Comment("twitter meta tags 'card'"),

		field.String("twitter_url").
			Default("").
			Comment("twitter meta tags 'url'"),

		field.String("twitter_title").
			Default("").
			Comment("twitter meta tags 'title'"),

		field.String("twitter_description").
			Default("").
			Comment("twitter meta tags 'description'"),

		field.String("twitter_image").
			Default("").
			Comment("twitter meta tags 'image'"),

		field.String("og_site_name").
			Default("").
			Comment("og meta tags 'site_name'"),

		field.String("og_locale").
			Default("").
			Comment("og meta tags 'locale'"),

		field.String("og_title").
			Default("").
			Comment("og meta tags 'title'"),

		field.String("og_description").
			Default("").
			Comment("og meta tags 'description'"),

		field.String("og_type").
			Default("").
			Comment("og meta tags 'type'"),

		field.String("og_url").
			Default("").
			Comment("og meta tags 'url'"),

		field.String("og_image").
			Default("").
			Comment("og meta tags 'image'"),

		field.String("og_image_type").
			Default("").
			Comment("og meta tags 'image:type'"),

		field.String("og_image_url").
			Default("").
			Comment("og meta tags 'image:url'"),

		field.String("og_image_secure_url").
			Default("").
			Comment("og meta tags 'image:secure_url'"),

		field.String("og_image_width").
			Default("").
			Comment("og meta tags 'image:width'"),

		field.String("og_image_height").
			Default("").
			Comment("og meta tags 'image:height'"),

		field.String("og_video").
			Default("").
			Comment("og meta tags 'video'"),

		field.String("og_video_type").
			Default("").
			Comment("og meta tags 'video:type'"),

		field.String("og_video_url").
			Default("").
			Comment("og meta tags 'video:url'"),

		field.String("og_video_secure_url").
			Default("").
			Comment("og meta tags 'video:secure_url'"),

		field.String("og_video_width").
			Default("").
			Comment("og meta tags 'video:width'"),

		field.String("og_video_height").
			Default("").
			Comment("og meta tags 'video:height'"),
	}
}

// Edges of the Page.
func (Page) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("page_source", PageSource.Type),
	}
}

// Indexes of the Page
func (Page) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("url").StorageKey("ux_url").Unique(),
		index.Fields("crawling_version").StorageKey("ux_url_and_crawling_version").Unique(),
	}
}

// Annotations of the Page.
func (Page) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "page",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_0900_ai_ci",
		},
	}
}
