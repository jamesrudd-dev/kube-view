import React from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import ClusterSelectionMenu from './components/ClusterSelectionMenu';

function App() {
  return (
    <div className="App">
      <h1 className="display-4 app-header">Kube View</h1>
      <p className="lead">A display of deployed images inside your kubernetes clusters and environments</p>
      <div>
        <ClusterSelectionMenu />
      </div>
    </div>
  );
}

export default App;
