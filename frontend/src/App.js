import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import { useGoogleLogin } from '@react-oauth/google';
// react-google-login DEPRECATED https://medium.com/@manwarerutuj/react-google-login-deprecated-dont-worry-react-oauth-google-will-fix-it-1b2141252952


function App() {
    const [user, setUser] = useState(null);
    const [profile, setProfile] = useState(null);

    const login = useGoogleLogin({
        onSuccess: (codeResponse) => setUser(codeResponse),
        onError: (error) => alert('Login nicht erfolgreich')
    });

    useEffect(() => {
        if (user) {
            axios.get(`https://www.googleapis.com/oauth2/v1/userinfo?access_token=${user.access_token}`, {
                headers: {
                    Authorization: `Bearer ${user.access_token}`,
                    Accept: 'application/json'
                }
            })
            .then((res) => setProfile(res.data))
            .catch((err) => alert(err));
        }
    }, [user]);

    function toggleTor(torNummer) {
      const gateName = `open${torNummer}`;
      const bodyData = {
        secret: "ultra-geheim", // TODO? Löschen?
        gate: gateName,
        //userEmail: user ? user.email : null,
        userEmail: "localtest@b-ite.de",
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
          <div><h1>Hallo, <br></br> ${user.email} </h1></div> // if user eingeloggt
        ) : (
          <span></span> // alternative? return no element oder so?
        )}
        <div>
        <img
              className="logo"
              src="https://www.b-ite.de/assets/images/bite-logo-9606eea416a94915.svg"
              alt="BITE GmbH Logo"
            />
        {!user ? (
          <div>
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
          <div className="button-container">
            <button className="neu-button" onClick={login}>Einloggen mit Google</button>
          </div>
        )}
        </div>
      </div>
    );
}

export default App;