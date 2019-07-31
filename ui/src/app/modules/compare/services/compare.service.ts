import { Injectable } from '@angular/core';

import { BACKEND_URL } from '../../../app.component'
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { CloudRegion } from '../components/compare-clouds/cloud-region';

@Injectable({
  providedIn: 'root'
}
)

export class CompareService{

    regions : CloudRegion[] = [];
    getRegionsUrl = BACKEND_URL + "clouds/regions";
    getCloudUrl = BACKEND_URL + "clouds/compare";
    postCloudRegion = BACKEND_URL + "clouds/send/region";

    constructor(private http: HttpClient){ }

    getRegions() : Observable<any>{
        return this.http.get<any>(this.getRegionsUrl);     
    }

    getCloudDetails() : Observable<any>{
        return this.http.get<any>(this.getCloudUrl);
    }

    sendCloudRegion():Observable<any>{
      return this.http.post<any>(this.postCloudRegion,null);
    }
    
}