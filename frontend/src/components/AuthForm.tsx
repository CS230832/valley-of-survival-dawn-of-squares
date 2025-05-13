import { createSignal } from "solid-js";
import { setSessionID } from "../utils/session";

interface Props {
    onAuthSuccess: () => void;
}

const AuthForm = (props: Props) => {
    const [username, setUsername] = createSignal("");
    const [password, setPassword] = createSignal("");
    const [isSignup, setIsSignup] = createSignal(false);
    const [error, setError] = createSignal<string | null>();
    const [message, setMessage] = createSignal<string | null>();
    const [signupSuccess, setSignupSuccess] = createSignal(false);

    const handleSubmit = async (e: Event) => {
        e.preventDefault();
        const endpoint = isSignup() ? "/api/signup" : "/api/login";

        try {
            const res = await fetch(endpoint, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    username: username(),
                    password: password(),
                }),
            });

            if (!res.ok) {
                setError(isSignup() ? "Signup failed" : "Login failed");
                setMessage(null);
                return;
            }

            setError(null);
            if (isSignup()) {
                setMessage("Signup successful! You can now log in.");
                setSignupSuccess(true);
            } else {
                const data = await res.text();
                setSessionID(data);
                props.onAuthSuccess();
            }
        } catch (err: Error | any) {
            setError(String(err));
            setMessage(null);
        }
    };

    const resetForm = () => {
        setUsername("");
        setPassword("");
        setIsSignup(false);
        setSignupSuccess(false);
        setError(null);
        setMessage(null);
    };

    const toggleSingupLogin = () => {
        setUsername("");
        setPassword("");
        setError(null);
        setMessage(null);
        setIsSignup(!isSignup());
    };

    return (
        <form class="auth-form" onSubmit={handleSubmit}>
            <h2 class="form-title">{isSignup() ? "Sign Up" : "Log In"}</h2>

            {signupSuccess() ? (
                <>
                    <p class="message">{message()}</p>
                    <button
                        class="submit-btn"
                        type="button"
                        onClick={resetForm}
                    >
                        Return to Login
                    </button>
                </>
            ) : (
                <>
                    <input
                        class="input-field"
                        placeholder="Username"
                        value={username()}
                        onInput={(e) => setUsername(e.currentTarget.value)}
                        required
                    />
                    <input
                        class="input-field"
                        type="password"
                        placeholder="Password"
                        value={password()}
                        onInput={(e) => setPassword(e.currentTarget.value)}
                        required
                    />
                    <button class="submit-btn" type="submit">
                        {isSignup() ? "Sign Up" : "Log In"}
                    </button>

                    {error() && <p class="error">{error()}</p>}
                    {message() && <p class="message">{message()}</p>}

                    <p>
                        {isSignup()
                            ? "Already have an account?"
                            : "Don't have an account?"}{" "}
                        <a href="#" onClick={() => toggleSingupLogin()}>
                            {isSignup() ? "Log in" : "Sign up"}
                        </a>
                    </p>
                </>
            )}
        </form>
    );
};

export default AuthForm;
