import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BACKEND_AUTH_URL } from '../../../app.constants';

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