import {NextPage} from "next";
import styles from "../styles/Home.module.css";
import {FC, useEffect, useState} from "react";
import React, {Component} from 'react';
import {
	AuthProviders,
	useLiveQuery,
	useMutation,
	useWunderGraph,
	withWunderGraph
} from "../components/generated/wundergraph.nextjs.integration";
import {AuthenticationClient} from 'authing-js-sdk';
import {
	BrowserRouter as Router,
	Switch,
	Route,
	useLocation,
	useHistory
} from "react-router-dom";

const authing = new AuthenticationClient({
	appId: '62a32077f02e938636734416',
	appHost: 'https://shaoxiong.authing.cn',
	redirectUri: 'http://localhost:3000/authentication',
	tokenEndPointAuthMethod: 'none'
});


const JobsPage: NextPage = () => {

	const {user, login, logout} = useWunderGraph();
	const [city, setCity] = useState<string>("Berlin");
	const onClick = () => {
		if (user === null || user === undefined) {
			login(AuthProviders.authing);
		} else {
			logout();
		}
	};
	return (
		<div className={styles.examplesContainer}>
			<h2>Login State</h2>
			<div>{user ? <p>{JSON.stringify(user)}</p> : <p>Not logged in</p>}</div>
			<button onClick={onClick}>{!user ? "Login" : "Logout"}</button>
			<input value={city} onChange={(e) => setCity(e.target.value)}/>
			<LoginBtn></LoginBtn>
			{/*<FindBlods />*/}
			{/*<InsertBlods />*/}
			{/* <App></App> */}
		</div>
	);
};

function LoginBtn() {
	return (
		<button
			onClick={() => {
				// PKCE 场景使用示例
				// 生成一个 code_verifier
				let codeChallenge = authing.generateCodeChallenge();
				localStorage.setItem('codeChallenge', codeChallenge);
				// 计算 code_verifier 的 SHA256 摘要
				let codeChallengeDigest = authing.getCodeChallengeDigest({codeChallenge, method: 'S256'});
				// 构造 OIDC 授权码 + PKCE 模式登录 URL
				let url = authing.buildAuthorizeUrl({codeChallenge: codeChallengeDigest, codeChallengeMethod: 'S256'});
				window.location.href = url;
			}}
		>
			登录
		</button>
	);
}

// function HandleCallback() {
//     let location = useLocation();
//     let query = new URLSearchParams(location.search);
//     let code = query.get('code');
//     let codeChallenge = localStorage.getItem('codeChallenge');
//     let history = useHistory();
//     useEffect(() => {
//         (async () => {
//             let tokenSet = await authing.getAccessTokenByCode(code, { codeVerifier: codeChallenge });
//             const { access_token, id_token } = tokenSet;
//             let userInfo = await authing.getUserInfoByAccessToken(tokenSet.access_token);
//             localStorage.setItem('accessToken', access_token);
//             localStorage.setItem('idToken', id_token);
//             localStorage.setItem('userInfo', JSON.stringify(userInfo));
//             history.push('/');
//         })();
//     });
//     return <div>加载中...</div>;
// }

// function HomePage() {
//     const [isAuthenticated, setIsAuthenticated] = useState()
//     useEffect(() => {
//         let userInfo = localStorage.getItem('userInfo');
//         setIsAuthenticated(!!userInfo);
//     }, []);
//     return (
//         <div className="App">
//             {isAuthenticated ? (<><LogoutBtn /></>) : <LoginBtn />}

//             <div>
//                 <div>用户信息：</div>
//                 <Profile />
//             </div>
//         </div>
//     )
// }

// function Profile() {
//     let userInfo = localStorage.getItem('userInfo');
//     return userInfo ?? '无';
// }


function LogoutBtn() {
	return (
		<button
			onClick={() => {
				let idToken = localStorage.getItem('idToken');
				localStorage.clear();
				let url = authing.buildLogoutUrl({expert: true, redirectUri: 'http://localhost:4000', idToken});
				window.location.href = url;
			}}
		>
			登出
		</button>
	);
}

// function App() {
//     return (
//         <Router>
//             <Switch>
//                 <Route path="/" exact={true}>
//                     <HomePage />
//                 </Route>
//                 <Route path="/callback" exact={true}>
//                     <HandleCallback />
//                 </Route>
//             </Switch>
//         </Router>
//     );
// }


// const InsertBlods: FC<{}> = () => {
//     const { result: r2, mutate: editMutate } = useMutation.CreateOneBlog()
//     const onSubmit = async (e: React.FormEvent<Element>) => {
//         e.preventDefault();
//         console.log(e)
//         const formIns = e.target;
//         editMutate({
//             input: { title: formIns[0].value, content: formIns[1].value, types: parseInt(formIns[2].value) }
//         })
//     }
//     return (
//         <form onSubmit={onSubmit}>
//             <input id="title" type="text" />
//             <input id="content" type="text" />
//             <input id="type" type="number" />
//             <button type="submit">添加</button>
//         </form>
//     )
// }

// const FindBlods: FC<{}> = () => {
//     //const findBlods = useLiveQuery.FindBlogs({});
//     const { result, mutate } = useMutation.RemoveBlodById()
//
//     const { result: r2, mutate: editMutate } = useMutation.EditBlogById()
//     const onSubmit = async (e: React.FormEvent<Element>) => {
//         e.preventDefault();
//         console.log(e)
//         const formIns = e.target;
//         editMutate({
//             input: { id: parseInt(formIns[0].value), title: formIns[1].value, content: formIns[2].value, types: parseInt(formIns[3].value) }
//         })
//     }
//     return (
//         <div>
//             <table>
//                 {findBlods.result.status === "ok" && (
//                     findBlods.result.data.FindBlogs.map((item) => {
//                         return (
//                             <ul>
//                                 <li>
//                                     <form onSubmit={onSubmit}>
//                                         <input id="id" type="number" value={item.id} />
//                                         <input id="title" type="text" defaultValue={item.title} />
//                                         <input id="content" type="text" defaultValue={item.content} />
//                                         <input id="type" type="number" defaultValue={item.types} />
//                                         <button type="submit">修改</button>
//                                         <button onClick={(e) => mutate({ input: { id: item?.id } })}>删除</button>
//                                     </form>
//                                 </li>
//                             </ul>
//                         )
//                     })
//                 )}
//             </table>
//         </div >
//     )
// }

const UpdateBlog: FC<{ title: string, content: string }> = ({title, content}) => {
	const onSubmit = async (e: React.FormEvent<Element>) => {
		e.preventDefault();
		const formData = new FormData();
	}
	return (
		<form onSubmit={onSubmit}>
			<input type="text" value={title}/>
			<input type="text" value={content}/>
			<button type="submit">Submit</button>
		</form>
	)
}


// const UpdateBlog: FC<{ id: number, title: string, content: string }> = ({ id, title, content }) => {
//     const { result, mutate } = useMutation.EditBlogById();
//     const liveWeather = mutate({
//         input: { id: id, title: title, content: content },
//     });
// }

// const ProtectedLiveWeather: FC<{ city: string }> = ({ city }) => {
//     const liveWeather = useLiveQuery.ProtectedWeather({
//         input: { forCity: city },
//     });
//     return (
//         <div>
//             {liveWeather.result.status === "requires_authentication" && (
//                 <p>LiveQuery waiting for user to be authenticated</p>
//             )}
//             {liveWeather.result.status === "loading" && <p>Loading...</p>}
//             {liveWeather.result.status === "error" && <p>Error</p>}
//             {liveWeather.result.status === "ok" && (
//                 <div>
//                     <h3>City: {liveWeather.result.data.getCityByName.name}</h3>
//                     <p>
//                         {JSON.stringify(liveWeather.result.data.getCityByName.coord)}
//                     </p>
//                     <h3>Temperature</h3>
//                     <p>
//                         {JSON.stringify(
//                             liveWeather.result.data.getCityByName.weather.temperature
//                         )}
//                     </p>
//                     <h3>Wind</h3>
//                     <p>
//                         {JSON.stringify(
//                             liveWeather.result.data.getCityByName.weather.wind
//                         )}
//                     </p>
//                 </div>
//             )}
//         </div>
//     );
// };

export default withWunderGraph(JobsPage);
