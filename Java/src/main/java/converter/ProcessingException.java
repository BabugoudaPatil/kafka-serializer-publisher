package converter;

public class ProcessingException extends RuntimeException {
    public ProcessingException(final String msg) {
        super(msg);
    }

    public ProcessingException(final String msg, Throwable cause) {
        super(msg, cause);
    }
}
