package science.monke.incoming;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.SneakyThrows;
import org.jdbi.v3.core.Jdbi;
import software.amazon.awssdk.services.secretsmanager.SecretsManagerClient;
import software.amazon.awssdk.services.secretsmanager.SecretsManagerClientBuilder;
import software.amazon.awssdk.services.secretsmanager.model.GetSecretValueRequest;
import software.amazon.awssdk.services.secretsmanager.model.GetSecretValueResponse;

import java.net.URI;

public class Utils {

  private Utils() {}

  @SneakyThrows
  public static Secret getSecret(final String secretId) {
    SecretsManagerClientBuilder secretsManagerClientBuilder = SecretsManagerClient.builder();

    String awsEndpointUrl = System.getenv("AWS_ENDPOINT_URL");
    if (awsEndpointUrl != null && !awsEndpointUrl.isEmpty()) {
      secretsManagerClientBuilder.endpointOverride(URI.create(awsEndpointUrl));
    }
    SecretsManagerClient secretsManagerClient = secretsManagerClientBuilder.build();

    GetSecretValueRequest getSecretValueRequest =
        GetSecretValueRequest.builder().secretId(secretId).build();
    GetSecretValueResponse getSecretValueResponse =
        secretsManagerClient.getSecretValue(getSecretValueRequest);
    secretsManagerClient.close();

    ObjectMapper objectMapper = new ObjectMapper();
    return objectMapper.readValue(getSecretValueResponse.secretString(), Secret.class);
  }

  public static Jdbi createDatabaseConnection(final DatabaseConfig databaseConfig) {
    final String url =
        String.format(
            "jdbc:postgresql://%s:%s/%s",
            databaseConfig.host, databaseConfig.port, databaseConfig.name);
    return Jdbi.create(url, databaseConfig.username, databaseConfig.password);
  }
}
