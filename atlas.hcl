// Define an environment named "local"
// (can also be unnamed `env`)
env {
  // Define where the schema definition resides.
  src = "file://sql/schema.sql"

  // Define the URL of the database for this environment.
  // (can also pass this in via environment variables).
  url = "postgres://postgres:admin@localhost:5432/flick"

  // Define the URL of the Dev Database for this environment
  // (used as a temp DB for schema validation during migrations)
  dev = "docker://postgres/17/dev"

  // Formatting options
  format {
    migrate {
      // Format to apply when running `atlas migrate diff <...>`
      diff = "{{ sql . \"\" }}"
    }
  }

  // Configure migrations
  migration {
    // URL where the migration directory resides.
    dir = "file://migrations"
    to = "file://sql/schema.sql"
    dev-url = "docker://postgres/17/dev"
  }
}