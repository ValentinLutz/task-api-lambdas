package science.monke.incoming;

public class DatabaseConfig {
  final String host;
  final String port;
  final String name;
  final String username;
  final String password;

  public DatabaseConfig(final Secret secret) {
    username = secret.username();
    password = secret.password();

    host = System.getenv("DB_HOST");
    port = System.getenv("DB_PORT");
    name = System.getenv("DB_NAME");
  }
}
