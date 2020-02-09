import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { tap, catchError } from 'rxjs/operators';
import { UrlData, ShortenData } from './app.model';

@Injectable({
  providedIn: 'root',
})
export class Service {

  constructor(
    private http: HttpClient) { }

  postUrl = '/api/post';
  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  };

  shorten(url: string): Observable<ShortenData> {
    const request: UrlData = { url };
    return this.http.post<ShortenData>(this.postUrl, request, this.httpOptions).pipe(
      catchError(this.handleError<ShortenData>('shorten'))
    );
  }


  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead


      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
