import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BACKEND_AUTH_URL } from '../../../app.component'

@Injectable()
export class LoginService {
    url: string;
    constructor(private http: HttpClient) {
        this.url = BACKEND_AUTH_URL + 'login';
    }
}