import Search from './Search';
import './App.css';
 

function App() {
  return (
    <div className="App">
      <header className="App-header" style={{"background-color":"#f9f9f9"}}>
        <h1> Emoji Search</h1>
        <Search />
      </header>
    </div>
  );
}

export default App;
