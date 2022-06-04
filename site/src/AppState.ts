class LoginDetails {
	username: string;
	token: string;
	userId: number;
	loggedIn: boolean;

    constructor(username: string, token: string, userId: number, loggedIn: boolean) {
        this.userId = userId;
        this.username = username;
        this.token = token;
        this.loggedIn = loggedIn;
    }
}

class AppProps {
	appState: AppState

    constructor(appState: AppState) {
        this.appState = appState;
    }
}

class AppState {
	loginState: LoginDetails | undefined;
	setLoginState: React.Dispatch<React.SetStateAction<LoginDetails | undefined>>;

	constructor(loginReactState: LoginDetails | undefined, setLoginReactState: React.Dispatch<React.SetStateAction<LoginDetails | undefined>>) {
		this.loginState = loginReactState;
		this.setLoginState = setLoginReactState;
	}
}

export { LoginDetails, AppState, AppProps };