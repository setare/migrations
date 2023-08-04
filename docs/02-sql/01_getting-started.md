---
sidebar_position: 1
---

# Getting Started

The `migrations-sql` package implements migrations based on the standard `database/sql`
package.

## What do migrations look like in my project?

For the SQL databases, we can write migrations as `.sql` files that will be stored
in the filesystem. So, in your project that can be a `migrations` directory:

```
./migrations
├── 20211015015442_create_customers_table.sql
├── 20211015015556_add_age_to_customers_table.sql
└── 20211015045556_add_birthday_to_customers_table.sql
```

By default, the `migrations` package will not enable the undoing of migrations. But, if you
enable it, you would find:

```
./migrations
├── 20211015015442_create_customers_table.do.sql
├── 20211015015442_create_customers_table.undo.sql
├── 20211015015556_add_age_to_customers_table.do.sql
├── 20211015015556_add_age_to_customers_table.undo.sql
└── 20211015045556_add_birthday_to_customers_table.do.sql
└── 20211015045556_add_birthday_to_customers_table.undo.sql
```

## How do I trigger the migrations on my service?

Create a file at `src/pages/my-markdown-page.md`:

```mdx title="src/pages/my-markdown-page.md"
# My Markdown page

This is a Markdown page
```

A new page is now available at `http://localhost:3000/my-markdown-page`.
