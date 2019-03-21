import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BACKEND_AUTH_URL } from '../../../app.component'

@Injectable()
export class LoginService {
    url: string;
    private sessionID: string;
    constructor(private http: HttpClient) {
        this.url = BACKEND_AUTH_URL + 'login';
    }

    public sendLoginCredential(credentials) {
        const httpPostOptions =
        {   
            withCredentials: true,
        }
        return this.http.post(this.url, credentials, httpPostOptions);
    }
}