import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BACKEND_AUTH_URL } from '../../../app.component'

@Injectable()
export class ChangepasswordService {
    url: string;
    constructor(private http: HttpClient) {
        this.url = BACKEND_AUTH_URL + 'changePassword';
    }

    public sendLoginCredentials(credentials) {
        return this.http.post(this.url, credentials);
    }
}