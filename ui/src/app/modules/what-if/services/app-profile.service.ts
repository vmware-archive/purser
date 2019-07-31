import { Injectable } from '@angular/core';
import { AppProfile } from '../components/app-profile';
import { BACKEND_URL } from '../../../app.component'
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
}
)
export class AppProfileService {
  appProfile : AppProfile;
  scenarioName : string;
  url_submit_spec : string;

  constructor(private http: HttpClient) {
    this.url_submit_spec = BACKEND_URL  + 'submit/whatif'
   }

  submitSpec(appProfile) : Observable<AppProfile>{
     return this.http.post<AppProfile>(this.url_submit_spec, appProfile);
   }
}
