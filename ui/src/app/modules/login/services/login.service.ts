import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BACKEND_AUTH_URL } from '../../../app.constants';

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