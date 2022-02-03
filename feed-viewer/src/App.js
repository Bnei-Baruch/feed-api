import React from 'react';
import {
  BrowserRouter as Router,
  Link,
  Route,
} from 'react-router-dom';

import './App.css';

import FeedContainer from './FeedContainer.js'
import Recommender from './Recommender.js'
import Chronicles from './Chronicles.js'

function App() {
  return (
  	<Router>
      <div className="App" style={{direction: 'rtl'}}>
        <Route exact path="/">
          <header className="App-header">
          <p>
            <span style={{textDecorationLine: 'underline'}}>Feed Viewer</span> / <Link to="/recommender">Recommender</Link> / <Link to="/chronicles">Chronicles</Link>
          </p>
          <div style={{width: '1200px'}}>
            <FeedContainer />
          </div>
          </header>
        </Route>
        <Route path="/recommender">
          <header className="App-header">
          <p>
            <Link to="/">Feed Viewer</Link> / <span style={{textDecorationLine: 'underline'}}>Recommender</span> / <Link to="/chronicles">Chronicles</Link>
          </p>
          <div style={{width: '1200px'}}>
            <Recommender />
          </div>
          </header>
        </Route>
        <Route path="/chronicles">
          <header className="App-header">
          <p>
            <Link to="/">Feed Viewer</Link> / <Link to="recommender">Recommender</Link> / <span style={{textDecorationLine: 'underline'}}>Chronicles</span>
          </p>
          <div style={{width: '100vw'}}>
            <Chronicles />
          </div>
          </header>
        </Route>
      </div>
  	</Router>
  );
}


export default App;
