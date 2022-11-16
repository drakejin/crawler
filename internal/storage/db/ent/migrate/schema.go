// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// PageColumns holds the columns for the "page" table.
	PageColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "referred_id", Type: field.TypeString},
		{Name: "crawling_version", Type: field.TypeString},
		{Name: "domain", Type: field.TypeString, Unique: true, Size: 700},
		{Name: "port", Type: field.TypeString, Size: 30, Default: "80"},
		{Name: "is_https", Type: field.TypeBool, Default: false},
		{Name: "indexed_url", Type: field.TypeString, Size: 1000, Default: ""},
		{Name: "path", Type: field.TypeString, Size: 1000, Default: ""},
		{Name: "querystring", Type: field.TypeString, Size: 2147483647, Default: "", SchemaType: map[string]string{"mysql": "text"}},
		{Name: "url", Type: field.TypeString, Size: 2147483647, Default: "", SchemaType: map[string]string{"mysql": "text"}},
		{Name: "count_referred", Type: field.TypeInt64, Default: 0},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"ALLOW", "NOTALLOW", "DELETE"}, Default: "ALLOW"},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP", SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "created_by", Type: field.TypeString, Size: 300},
		{Name: "updated_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP", SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_by", Type: field.TypeString, Size: 300},
		{Name: "title", Type: field.TypeString, Default: ""},
		{Name: "description", Type: field.TypeString, Default: ""},
		{Name: "keywords", Type: field.TypeString, Default: ""},
		{Name: "content_language", Type: field.TypeString, Default: ""},
		{Name: "twitter_card", Type: field.TypeString, Default: ""},
		{Name: "twitter_url", Type: field.TypeString, Default: ""},
		{Name: "twitter_title", Type: field.TypeString, Default: ""},
		{Name: "twitter_description", Type: field.TypeString, Default: ""},
		{Name: "twitter_image", Type: field.TypeString, Default: ""},
		{Name: "og_site_name", Type: field.TypeString, Default: ""},
		{Name: "og_locale", Type: field.TypeString, Default: ""},
		{Name: "og_title", Type: field.TypeString, Default: ""},
		{Name: "og_description", Type: field.TypeString, Default: ""},
		{Name: "og_type", Type: field.TypeString, Default: ""},
		{Name: "og_url", Type: field.TypeString, Default: ""},
		{Name: "og_image", Type: field.TypeString, Default: ""},
		{Name: "og_image_type", Type: field.TypeString, Default: ""},
		{Name: "og_image_url", Type: field.TypeString, Default: ""},
		{Name: "og_image_secure_url", Type: field.TypeString, Default: ""},
		{Name: "og_image_width", Type: field.TypeString, Default: ""},
		{Name: "og_image_height", Type: field.TypeString, Default: ""},
		{Name: "og_video", Type: field.TypeString, Default: ""},
		{Name: "og_video_type", Type: field.TypeString, Default: ""},
		{Name: "og_video_url", Type: field.TypeString, Default: ""},
		{Name: "og_video_secure_url", Type: field.TypeString, Default: ""},
		{Name: "og_video_width", Type: field.TypeString, Default: ""},
		{Name: "og_video_height", Type: field.TypeString, Default: ""},
	}
	// PageTable holds the schema information for the "page" table.
	PageTable = &schema.Table{
		Name:       "page",
		Columns:    PageColumns,
		PrimaryKey: []*schema.Column{PageColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "page_crawling_version",
				Unique:  false,
				Columns: []*schema.Column{PageColumns[2]},
			},
			{
				Name:    "page_referred_id",
				Unique:  false,
				Columns: []*schema.Column{PageColumns[1]},
			},
		},
	}
	// PageLinksColumns holds the columns for the "page_links" table.
	PageLinksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "source_id", Type: field.TypeInt64},
		{Name: "target_id", Type: field.TypeInt64},
		{Name: "created_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP", SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "created_by", Type: field.TypeString, Size: 300},
		{Name: "updated_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP", SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_by", Type: field.TypeString, Size: 300},
	}
	// PageLinksTable holds the schema information for the "page_links" table.
	PageLinksTable = &schema.Table{
		Name:       "page_links",
		Columns:    PageLinksColumns,
		PrimaryKey: []*schema.Column{PageLinksColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "pagelink_source_id_target_id",
				Unique:  true,
				Columns: []*schema.Column{PageLinksColumns[1], PageLinksColumns[2]},
			},
			{
				Name:    "pagelink_source_id",
				Unique:  false,
				Columns: []*schema.Column{PageLinksColumns[1]},
			},
			{
				Name:    "pagelink_target_id",
				Unique:  false,
				Columns: []*schema.Column{PageLinksColumns[2]},
			},
		},
	}
	// PageSourceColumns holds the columns for the "page_source" table.
	PageSourceColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "referred_page_id", Type: field.TypeString},
		{Name: "url", Type: field.TypeString, Size: 2147483647, Default: "", SchemaType: map[string]string{"mysql": "text"}},
		{Name: "referred_url", Type: field.TypeString, Size: 2147483647, Default: "", SchemaType: map[string]string{"mysql": "text"}},
		{Name: "source", Type: field.TypeString, Size: 2147483647, Default: "", SchemaType: map[string]string{"mysql": "text"}},
		{Name: "page_id", Type: field.TypeString},
	}
	// PageSourceTable holds the schema information for the "page_source" table.
	PageSourceTable = &schema.Table{
		Name:       "page_source",
		Columns:    PageSourceColumns,
		PrimaryKey: []*schema.Column{PageSourceColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "page_source_page_page_source",
				Columns:    []*schema.Column{PageSourceColumns[5]},
				RefColumns: []*schema.Column{PageColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		PageTable,
		PageLinksTable,
		PageSourceTable,
	}
)

func init() {
	PageTable.Annotation = &entsql.Annotation{
		Table:     "page",
		Charset:   "utf8mb4",
		Collation: "utf8mb4_0900_ai_ci",
	}
	PageLinksTable.Annotation = &entsql.Annotation{
		Table:     "page_links",
		Charset:   "utf8mb4",
		Collation: "utf8mb4_0900_ai_ci",
	}
	PageSourceTable.ForeignKeys[0].RefTable = PageTable
	PageSourceTable.Annotation = &entsql.Annotation{
		Table:     "page_source",
		Charset:   "utf8mb4",
		Collation: "utf8mb4_0900_ai_ci",
	}
}
