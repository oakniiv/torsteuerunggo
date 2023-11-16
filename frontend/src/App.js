import React, { useState } from 'react';
import { GoogleLogin } from 'react-google-login';
import './App.css';

function App() {
  const [user, setUser] = useState(null);
  const responseGoogle = (response) => {
    if (response.profileObj) {
      setUser({
        email: response.profileObj.email,
      });
    }
  };
  function toggleTor(torNummer) {
    const gateName = `open${torNummer}`; //war früher value in der Form
    const bodyData = {
      secret: "ultra-geheim", //TODO? löschen?
      gate: gateName,
      userEmail: user ? user.email : null,
    };
    fetch('http://127.0.0.1:8080/api/toggle', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(bodyData),
    })
      .then(function (response) {
        if (!response.ok) {
          throw new Error('error');
        }
      })
      .catch(function (error) {
        console.error('Error:', error);
      });
  }
  return (
    <div className="container">
      {user ? (
        <div><h1>Hallo, {user.email}</h1></div>
      ) : (
        <GoogleLogin
          clientId="client id????"
          buttonText="Google Login"
          onSuccess={responseGoogle}
          onFailure={responseGoogle}
          cookiePolicy={'single_host_origin'}
          isSignedIn={true} // isSignedIn={true} attribute will call onSuccess callback on load to keep the user signed in.
        />
        // https://www.npmjs.com/package/react-google-login
      )}
      <br></br>
      <br></br>
      {user ? (
        <div>
          <img
            className="logo"
            src="https://www.b-ite.de/assets/images/bite-logo-9606eea416a94915.svg"
            alt="BITE GmbH Logo"
          />
          <div className="button-container">
            {[1, 2, 3, 4].map((torNummer) => (
              <button
                key={torNummer}
                className="neu-button"
                onClick={() => toggleTor(torNummer)}
              >
                <strong>Tor {torNummer}</strong>
              </button>
            ))}
          </div>
        </div>
      ) : (
        <div>
          <img
            className="logo"
            src="https://www.b-ite.de/assets/images/bite-logo-9606eea416a94915.svg"
            alt="BITE GmbH Logo"
          />
          Du musst dich einloggen.
        </div>
        // https://www.npmjs.com/package/react-google-login
      )}
    </div>
  );
}

export default App;
