import { Component } from '@angular/core';
import { Service } from './app.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'url-shortener';
  shorten = '';

  constructor(private service: Service) {
  }

  onClick(url: string) {
    this.service.shorten(url)
      .subscribe(response => this.shorten = response.shortenUrl);
  }
}
