import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';

import './App.css';

import FeedContainer from './FeedContainer.js'

function App() {
  return (
  	<Router>
		<div className="App" style={{direction: 'rtl'}}>
		  <header className="App-header">
			<p>
			  Feed Viewer
			</p>
			<div style={{width: '800px'}}>
				<FeedContainer/>
			</div>
		  </header>
		</div>
  	</Router>
  );
}

export default App;
