import java.io.*;
import java.net.*;
import java.nio.charset.*;
import com.google.gson.*;

// This is an executable example of how to call the import-bibtex service from Java.

public abstract class ImportBibtex {
    public static void main(String[] args) {
        try {
            // The URL of the service. Should be on localhost with access to a shared database.
            URL url = new URL("http://localhost:8080/import-bibtex");
            HttpURLConnection conn = (HttpURLConnection)url.openConnection();
            // Some precautions.
            conn.setConnectTimeout(5000);
            conn.setReadTimeout(5000);
            // Send the BibTeX data. In production the data shall come from a web form.
            conn.setDoOutput(true);
            conn.setRequestMethod("POST");
            OutputStreamWriter writer = new OutputStreamWriter(conn.getOutputStream(), StandardCharsets.UTF_8);
            writer.write("""
            @article{an-id,
                title = "A Title",
                author = "John Smith and Mary Stone"
            }
            """);
            writer.close();
            // Check to response status.
            System.out.println("response status: " + conn.getResponseCode() + " (" + conn.getResponseMessage() + ")");
            // Read the returned JSON data.
            if (conn.getResponseCode() == 200) {
                InputStreamReader reader = new InputStreamReader(conn.getInputStream(), StandardCharsets.UTF_8);
                JsonObject obj = new JsonParser().parse(reader).getAsJsonObject();
                System.out.println(obj);
            }
        } catch (Exception ex) {
            System.err.println("failed to connect to service: " + ex);
        }
    }
}
