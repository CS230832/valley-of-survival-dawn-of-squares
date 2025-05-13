import "./App.css";
import { createSignal, onMount, Show } from "solid-js";
import { getSessionID, removeSessionID } from "./utils/session";
import { apiFetch } from "./api/client";
import AuthForm from "./components/AuthForm";
import GameMenu from "./components/GameMenu";
import Spawn from "./components/Spawn";
import ClanMenu from "./components/ClanMenu";

function App() {
    const [loggedIn, setLoggedIn] = createSignal(false);
    const [spawned, setSpawned] = createSignal(false);
    const [inManageClan, setInManageClan] = createSignal(false);

    const logout = () => {
        removeSessionID();
        setLoggedIn(false);
        setSpawned(false);
    };

    onMount(async () => {
        const session = getSessionID();
        if (session) {
            try {
                const res = await apiFetch("/api/verify_session");
                if (!res.ok) throw new Error("Invalid session");
                setLoggedIn(true);
            } catch {
                logout();
            }
        }
    });

    return (
        <main class="main-container">
            <Show
                when={loggedIn()}
                fallback={<AuthForm onAuthSuccess={() => setLoggedIn(true)} />}
            >
                <Show
                    when={spawned()}
                    fallback={
                        <Show
                            when={inManageClan()}
                            fallback={
                                <GameMenu
                                    onLogout={logout}
                                    onManageClan={() => setInManageClan(true)}
                                    onSpawn={() => setSpawned(true)}
                                />
                            }
                        >
                            <ClanMenu
                                onReturn={() => setInManageClan(false)}
                                onLogout={logout}
                            />
                        </Show>
                    }
                >
                    <Spawn />
                </Show>
            </Show>
        </main>
    );
}

export default App;
