package science.monke.incoming;

import com.fasterxml.jackson.annotation.JsonProperty;

public record TaskRequest(
    @JsonProperty("title") String title, @JsonProperty("description") String description) {}
