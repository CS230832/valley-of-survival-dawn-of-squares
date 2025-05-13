interface Props {
    onLogout: () => void;
    onManageClan: () => void;
    onSpawn: () => void;
}

const GameMenu = (props: Props) => {
    return (
        <div class="game-menu">
            <h2>Game Menu</h2>
            <button onClick={props.onSpawn}>Spawn into Game</button>
            <button onClick={props.onManageClan}>Manage Clan</button>
            <button onClick={props.onLogout}>Logout</button>
        </div>
    );
};

export default GameMenu;
